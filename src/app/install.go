package app

import (
	"context"
	"errors"
	"fmt"
	"horus/app/admin/command"
	"horus/src"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/gta-spec/utils/file"
)

// isInstall 判断是否安装过了
func isInstall() bool {
	filename := src.InstallPath + "install.lock"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// install 首次进入,启动系统安装
func install(engine *gin.Engine, address string) {
	// 加载模板
	engine.LoadHTMLFiles(src.InstallPath + "install.html")
	// 注册安装路由
	engine.Match([]string{"GET", "POST"}, "/install", new(command.Install).Index)

	engine.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/install")
	})
	// 创建HTTP服务器实例
	srv := &http.Server{
		Addr:    address,
		Handler: engine,
	}

	// 监听文件创建
	complete := make(chan fsnotify.Event)
	err := _file.Watcher(src.InstallPath+"install.lock", complete, fsnotify.Create|fsnotify.Rename)

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
