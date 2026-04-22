package app

import (
	"context"
	"errors"
	"fmt"
	"gota/internal/admin/command"
	"gota/pkg"
	"gota/pkg/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

// isInstall 判断是否安装过了
func isInstall() bool {
	filename := pkg.INSTALL_PATH + "install.lock"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// install 首次进入,启动系统安装
func install(engine *gin.Engine, address string) {
	// 加载模板
	engine.LoadHTMLFiles(pkg.INSTALL_PATH + "install.html")
	// 注册安装路由
	engine.Any("/install", command.Install{MinGoVersion: "1.23"}.Index)

	engine.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/install")
	})

	// 创建完成channel
	complete := make(chan struct{})

	// 创建HTTP服务器实例
	srv := &http.Server{
		Addr:    address,
		Handler: engine,
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
