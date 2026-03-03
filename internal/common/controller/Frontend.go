package controller

import (
	"fmt"
	Config "gonet/pkg/config"
	"gonet/pkg/template"
	"gonet/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Frontend struct {
}

func (c Frontend) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {
		tpl := template.NewTemplate(c)
		siteViper, _ := Config.Get("site")
		uploadViper, _ := Config.Get("upload")

		site := map[string]any{}
		for _, v := range []string{"name", "cdnurl", "version", "timezone", "languages"} {
			switch v {
			case "languages":
				site[v] = siteViper.GetStringMapString(v)
			default:
				site[v] = siteViper.GetString(v)
			}
		}
		fmt.Println("我是Initialize")

		config := map[string]any{
			"app_debug":      gin.IsDebugging(),
			"site":           site,
			"upload":         uploadViper.AllSettings(),
			"modulename":     c.GetString("modulename"),
			"controllername": c.GetString("controllername"),
			"actionname":     c.GetString("actionname"),
			"jsname":         fmt.Sprintf("frontend/%s", c.GetString("controllername")),
			"moduleurl":      fmt.Sprintf("/%s", c.GetString("modulename")),
			"language":       utils.Langset(c.Request.Header.Get("Accept-Language")),
		}

		tpl.Assign("site", siteViper.AllSettings())
		tpl.Assign("config", config)

		c.Set("Think", tpl)
	}
}
