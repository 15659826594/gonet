package main

import (
	_ "gonet/internal/index/controller"
	"gonet/pkg"
	"gonet/pkg/app"
	_ "gonet/pkg/app/autoload"
	"gonet/pkg/config"
	_ "gonet/pkg/i18n"
	"gonet/pkg/logger"
	"gonet/pkg/logger/zap"
	"strings"
)

// @title OpenAPI
// @version 1.0
// @description Gonet框架
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
	// 加载配置
	_ = config.SetGlobalConfigFile("pkg/config/config.yaml")
	cfgInst := config.Viper()
	// 加载次要配置
	config.LoadConfigGlob(strings.Join([]string{pkg.APP_PATH, "extra", "*.yaml"}, "/"))
	// 初始化zap日志
	logInst := logger.ReplaceGlobals(zap.Adapter(zap.New(cfgInst.Log)))
	defer logInst.Sync()

	application := app.New(cfgInst)

	application.Run()
}
