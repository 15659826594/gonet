package app

import (
	"fmt"

	"horus/app/index/controller"

	"github.com/gin-gonic/gin"
)

var (
	Modules        = map[string]*gin.RouterGroup{}
	RegisterRoutes func()
)

func Router(engine *gin.Engine) {
	admin := engine.Group("admin")
	{
		fmt.Println(admin)
	}

	index := engine.Group("")

	{
		demo := index.Group("")

		{
			demo1 := new(controller.Index)
			demo.GET("", demo1.Index)
		}
	}

	api := engine.Group("/api")

	{
		demo := api.Group("demo")

		{
			fmt.Println(demo)
		}
	}
}
