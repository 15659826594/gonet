//go:generate go run src/cmd/tools/gen_route.go
//go:generate swag init -o ./src/docs

package main

import (
	"fmt"
	_ "horus/app/api/controller"
	_ "horus/app/index/controller"
	"horus/src"
	"horus/src/app"
	_ "horus/src/app/autoload"
	"horus/src/cmd/router"
	"horus/src/config"
	"horus/src/database"
	"horus/src/database/redis"
	_ "horus/src/i18n"
	"horus/src/logger"
	"strings"
)

// @title OpenAPI
// @version 1.0
// @description gt框架
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
	//mainPath := pkg.MainPath
	//outFilename := filepath.Join("./", root, "route_gen.go")
	//out, err := os.Create(outFilename)
	//if err != nil {
	//	panic(err)
	//}
	//defer out.Close()
	s, _ := router.ScanPackages("app")
	for _, p := range s {
		//fmt.Println(p.Path)
		ast := router.NewFileAst(p.Path)
		for _, i2 := range p.Files {
			_ = ast.ParseFile(i2)
		}
		for _, structType := range ast.StructTypes {
			fmt.Println(structType.Docs)
		}
		//fmt.Println("*************")
	}
	return
	//out.WriteString(`package app`)
	//if mainPath != "" {
	//	return
	//}

	defer logger.Close()
	defer database.Close()
	defer redis.Close()
	// 加载配置
	_ = config.SetGlobalConfigFile("src/config/config.yaml")
	cfgInst := config.Viper()
	// 加载次要配置
	config.LoadConfigGlob(strings.Join([]string{src.AppPath, "extra", "*.yaml"}, "/"))

	application := app.New(cfgInst)

	application.Run()
}
