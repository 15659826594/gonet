package main

import (
	_ "gota/internal/api/controller"
	_ "gota/internal/index/controller"
	"gota/pkg"
	"gota/pkg/app"
	_ "gota/pkg/app/autoload"
	"gota/pkg/config"
	"gota/pkg/database"
	"gota/pkg/database/redis"
	_ "gota/pkg/i18n"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// @title OpenAPI
// @version 1.0
// @description Gota框架
// @termsOfService https://www.swagger.io/terms/

// @host 127.0.0.1:8080
// @BasePath /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://github.com/swaggo/swag/blob/master/README_zh-CN.md

// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API 认证方式：在请求头中添加 Authorization 字段，值为 Bearer + 空格 + token
func main() {
	os.Setenv("APP_ENV", gin.ReleaseMode)
	defer database.Close()
	defer redis.Close()
	// 加载配置
	_ = config.SetGlobalConfigFile("pkg/config/config.yaml")
	cfgInst := config.Viper()
	// 加载次要配置
	config.LoadConfigGlob(strings.Join([]string{pkg.APP_PATH, "extra", "*.yaml"}, "/"))

	application := app.New(cfgInst)

	application.Run()
}
