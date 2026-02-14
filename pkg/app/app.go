package app

import (
	"context"
	"errors"
	"fmt"
	_ "gonet/docs"
	"gonet/internal/admin/command"
	"gonet/pkg"
	"gonet/pkg/app/route"
	"gonet/pkg/config"
	"gonet/pkg/database"
	"gonet/pkg/i18n"
	"gonet/pkg/logger"
	"gonet/pkg/middleware"
	"gonet/pkg/template/multi"
	"gonet/pkg/utils"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var funcMap = template.FuncMap{
	"json_encode": utils.JsonEncode,
	"json_decode": utils.JsonDecode,
	"url":         new(utils.Url).Build,
	"cdnurl": func(url string) string {
		return url
	},
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type App struct {
	*gin.Engine
	Config   *config.Config
	MinGoVer string
}

func New(c *config.Config) *App {

	gin.SetMode(c.Mode())

	app := &App{
		Engine:   gin.Default(),
		Config:   c,
		MinGoVer: "1.23",
	}

	// 判断是否安装了项目
	if !IsInstall() {
		app.install()
		//重新加载配置
		//viper.Reset()
		//autoload.Load()
		//_ = config.SetGlobalConfigFile(pkg.CONF_PATH + "config.yaml")
		//app.Config = config.Viper()
	}

	app.Engine.Use(app.Config.ExceptionHandle(), middleware.ResponseHandler(), middleware.CorsMiddleware())

	dbViper, err := config.SetConfigFile("database", "pkg/database/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load database configuration: %v", err)
	}

	// 初始化gorm
	database.Gorm(new(database.Database).Viper(dbViper))

	//设置html模板
	HTMLRender := multi.Render(
		gin.IsDebugging(),
		multi.WithPrefix("[GIN-debug] "),
		multi.WithReplaces(app.Config.ViewReplaceStr),
		multi.WithDelims(app.Config.Template.Delims()),
		multi.WithFuncMap(funcMap),
	)
	HTMLRender.SetScopeFuncMap(scopeFuncMap)
	HTMLRender.LoadHTMLGlob("**/*.html")
	app.Engine.HTMLRender = HTMLRender

	return app
}

func (a *App) Run(addr ...string) {
	_, filename, line, _ := runtime.Caller(0)
	caller := filepath.ToSlash(strings.TrimPrefix(filepath.FromSlash(filename), filepath.FromSlash(pkg.RootPath()+string(filepath.Separator))))
	//设置路由
	route.Build(a.Engine, func(name string) (*gin.RouterGroup, string) {
		modulename := strings.Split(name, pkg.DS)[1]
		return a.Engine.Group(modulename), modulename
	})
	// 初始化静态资源
	a.Engine.Static("/assets", "./public/assets")
	a.Engine.StaticFile("/favicon.ico", "./assets/favicon.ico")
	a.Engine.GET("swagger.json", func(c *gin.Context) {
		filePath := "./docs/swagger.json"
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Swagger file not found"})
			return
		}
		c.File(filePath)
	})

	address := resolveAddress(addr)
	host := strings.Split(address, ":")

	_ = a.Engine.SetTrustedProxies([]string{host[0]})

	fmt.Println(fmt.Sprintf(`%s server running for the %s:%d process at:

	➜  Local:   http://%s/
	➜  Docs:    http://%s/swagger.json

start gin %s...`, gin.Mode(), caller, line-1, address, address, a.Config.AppNamespace))

	logger.Info(fmt.Sprintf("HTTP Server listening at %s", host[1]))

	if db, err := database.Gorm().DB(); err == nil {
		defer db.Close()
	}

	err := a.Engine.Run(address)
	if err != nil {
		log.Fatalf("gin run failed: %v", err)
	}
}

// IsInstall 判断是否安装过了
func IsInstall() bool {
	filename := pkg.INSTALL_PATH + "install.lock"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// install 首次进入,启动系统安装
func (a *App) install() {
	//a.LoadHTMLFolder(pkg.INSTALL_PATH, "*.html")

	action := command.Install{
		MinGoVersion: a.MinGoVer,
	}.Index

	// 注册安装路由
	a.Engine.GET("/install", action)
	a.Engine.POST("/install", action)

	a.Engine.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/install")
	})

	// 创建完成channel
	complete := make(chan struct{})

	// 创建HTTP服务器实例
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("APP_HOSTNAME"), viper.GetInt("APP_HOSTPORT")),
		Handler: a.Engine,
	}

	// 监听 install.lock 文件的创建
	err := utils.FileListener(pkg.INSTALL_PATH, func(event fsnotify.Event, done func()) {
		if (event.Op&fsnotify.Create == fsnotify.Create) && filepath.Base(event.Name) == "install.lock" {
			done()
			close(complete)
		}
	})

	if err != nil {
		log.Fatalf("File watcher error: %v", err)
	}

	// 阻塞主线程直到安装完成
	fmt.Printf("Please visit http://%s/install to complete installation\n", srv.Addr)

	// 启动临时服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// 等待安装完成
	<-complete

	// 安装完成后关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

// 子模板 scopeFuncMap 独立作用域
func scopeFuncMap(name string, fm template.FuncMap) (string, template.FuncMap) {
	links := strings.Split(name, "/")
	if len(links) > 1 && links[1] == "view" {
		links[1] = ""
	}

	link := filepath.ToSlash(filepath.Clean(strings.Join(links, "/")))
	link = strings.TrimSuffix(link, filepath.Ext(link))

	fm["__"] = func(messageID string, templates ...map[string]any) string {
		return i18n.T(link, messageID, templates...)
	}
	return name, fm
}

func resolveAddress(addr []string) string {
	host := viper.GetString("APP_HOSTNAME")
	if host == "" {
		host = "localhost"
	}
	switch len(addr) {
	case 0:
		if port := viper.GetInt("APP_HOSTPORT"); port != 0 {
			return fmt.Sprintf("%s:%d", host, port)
		}
		return fmt.Sprintf("%s:8080", host)
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
