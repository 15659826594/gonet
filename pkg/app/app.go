package app

import (
	"fmt"
	_ "gota/docs"
	"gota/pkg"
	"gota/pkg/app/route"
	"gota/pkg/config"
	"gota/pkg/logger"
	"gota/pkg/middleware"
	"gota/pkg/template/multi"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

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
		Engine: gin.New(),
		Config: c,
	}
	return app
}

func (t *App) Run(addr ...string) {
	_, filename, line, _ := runtime.Caller(0)
	caller := filepath.ToSlash(strings.TrimPrefix(filepath.FromSlash(filename), filepath.FromSlash(pkg.RootPath()+string(filepath.Separator))))

	address := resolveAddress(addr)
	host := strings.Split(address, ":")
	// 判断是否安装了项目
	if !isInstall() {
		install(t.Engine, address)
	}

	t.Engine.Use(middleware.ResponseHandler(), logger.Middleware(), middleware.Session(), middleware.Cors())

	//设置html模板
	rs := map[string]string{}
	for k, v := range t.Config.ViewReplaceStr {
		rs[strings.ToUpper(k)] = v
	}
	multi.SetReplaces(rs)
	multi.SetDelims(t.Config.Template.Delims())
	multi.SetFuncMap(funcMap)

	render := multi.New(gin.IsDebugging())

	render.LoadHTMLGlob(os.DirFS("./internal"), "**/*.{html,tpl}")
	render.LoadHTMLFile(t.Config.DispatchSuccessTmpl)
	render.LoadHTMLFile(t.Config.DispatchErrorTmpl)

	for name, tpl := range render.Seq2() {
		name = parseTplRoute(name)
		tpl.SetFuncMap("__", I18n(name))
		tpl.SetFuncMap("url", Url(name))
	}

	t.Engine.HTMLRender = render.Init()

	//设置路由
	router(t.Engine)

	_ = t.Engine.SetTrustedProxies([]string{host[0]})

	fmt.Println(fmt.Sprintf(`%s server running for the %s:%d process at:

	➜  Local:   http://%s/
	➜  Docs:    http://%s/swagger.json

start gin %s...`, gin.Mode(), caller, line-1, address, address, t.Config.AppNamespace))

	logger.Record(nil, slog.LevelInfo, fmt.Sprintf("HTTP Server listening at %s", host[1]))

	t.Engine.Run(address)
}
func router(engine *gin.Engine) {
	// 初始化静态资源
	engine.Static("/assets", "./public/assets")
	engine.StaticFile("/favicon.ico", "./assets/favicon.ico")
	route.Build(engine, func(name string) (*gin.RouterGroup, string) {
		modulename := strings.Split(name, pkg.DS)[1]
		return engine.Group(modulename), modulename
	})
	engine.GET("swagger.json", func(c *gin.Context) {
		filePath := "./docs/swagger.json"
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Swagger file not found"})
			return
		}
		c.File(filePath)
	})
}

func parseTplRoute(path string) string {
	re := regexp.MustCompile(`^([^/]+)/view(/.*)?$`)
	result := re.ReplaceAllString(path, "$1$2")
	return strings.TrimSuffix(result, filepath.Ext(result))
}

// 入口地址
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
