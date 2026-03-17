package controller

import (
	"fmt"
	"gonet/internal/common/library/Auth"
	"gonet/internal/common/model"
	Config "gonet/pkg/config"
	"gonet/pkg/library/think"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Frontend struct {
	think.Controller
}

func (t Frontend) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {
		modulename := c.GetString("modulename")
		controllername := c.GetString("controllername")
		actionname := c.GetString("actionname")

		//siteViper, _ := Config.Get("site")
		//uploadViper, _ := Config.Get("upload")
		//
		//site := map[string]any{}
		//for _, v := range []string{"name", "cdnurl", "version", "timezone", "languages"} {
		//	switch v {
		//	case "languages":
		//		site[v] = siteViper.GetStringMapString(v)
		//	default:
		//		site[v] = siteViper.GetString(v)
		//	}
		//}
		//
		auth := Auth.Instance()
		c.Set("auth", auth)

		// token
		token := token(c)

		path := strings.ReplaceAll(controllername, ".", "/") + "/" + actionname

		// 设置当前请求的URI
		auth.SetRequestUri(path)
		// 检测是否需要验证登录
		if !auth.Match(c.GetStringSlice("noNeedLogin")) {
			//初始化
			auth.Init(token)
			//检测是否登录
			if auth.IsLogin() {
				t.Error(c, "Please login first", "index/user/login")
			}
			// 判断是否需要验证权限
			if !auth.Match(c.GetStringSlice("noNeedRight")) {
				// 判断控制器和方法判断是否有对应权限
				if auth.Check(path, modulename) {
					t.Error(c, "You have no permission")
				}
			}
		} else {
			// 如果有传递token才验证是否登录状态
			if token != "" {
				auth.Init(token)
			}
		}

		t.View.Assign(c, "user", auth.GetUser())

		// 语言检测
		lang := langset(c).String()
		if matched, _ := regexp.MatchString(`^([a-zA-Z\-_]{2,10})$`, lang); !matched {
			lang = "zh-cn"
		}

		site, _ := Config.Get("site")

		upload := new(model.Config).Upload()

		config := map[string]any{
			"app_debug":      gin.IsDebugging(),
			"site":           mapIntersectKey(site.AllSettings(), sliceFlip([]string{"name", "cdnurl", "version", "timezone", "languages"})),
			"upload":         upload,
			"modulename":     modulename,
			"controllername": controllername,
			"actionname":     actionname,
			"jsname":         fmt.Sprintf("frontend/%s", controllername),
			"moduleurl":      fmt.Sprintf("/%s", modulename),
			"language":       lang,
		}

		t.View.Assign(c, "site", site.AllSettings())
		t.View.Assign(c, "config", config)
	}
}

func sliceFlip[T comparable](s []T) map[T]int {
	result := make(map[T]int, len(s))
	for i, v := range s {
		result[v] = i
	}
	return result
}
func mapIntersectKey[K comparable, V any, W any](first map[K]V, others ...map[K]W) map[K]V {
	if len(first) == 0 || len(others) == 0 {
		return make(map[K]V)
	}

	result := make(map[K]V)

	for key := range first {
		foundInAll := true

		for _, otherMap := range others {
			if _, exists := otherMap[key]; !exists {
				foundInAll = false
				break
			}
		}

		if foundInAll {
			result[key] = first[key]
		}
	}

	return result
}
