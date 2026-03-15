package controller

import (
	"gonet/pkg/library/think"

	"github.com/gin-gonic/gin"
)

type Frontend struct {
	think.Controller
}

func (t Frontend) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {

		//modulename := c.GetString("modulename")
		//controllername := c.GetString("controllername")
		//actionname := c.GetString("actionname")
		//
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
		//auth := Auth.Instance()
		//c.Set("auth", auth)
		//
		//// getToken
		//token := "getToken"
		//
		//path := strings.ReplaceAll(controllername, ".", "/") + "/" + actionname
		//
		////t.Assign("user", auth.GetUser())
		//
		//t.Error(c, "Please login first", "index/user/login")
		//
		//// 设置当前请求的URI
		//auth.SetRequestUri(path)
		//// 检测是否需要验证登录
		//if !auth.Match(c.GetStringSlice("noNeedLogin")) {
		//	//初始化
		//	auth.Init(token)
		//	//检测是否登录
		//	if auth.IsLogin() {
		//		t.Error(c, "Please login first", "index/user/login")
		//	}
		//	// 判断是否需要验证权限
		//	if !auth.Match(c.GetStringSlice("noNeedRight")) {
		//		// 判断控制器和方法判断是否有对应权限
		//		if auth.Check(path, modulename) {
		//			t.Error(c, "You have no permission")
		//		}
		//	}
		//} else {
		//	// 如果有传递token才验证是否登录状态
		//	if token != "" {
		//		auth.Init(token)
		//	}
		//}
		//
		//view.Assign("user", auth.GetUser())
		//
		//config := map[string]any{
		//	"app_debug":      gin.IsDebugging(),
		//	"site":           site,
		//	"upload":         uploadViper.AllSettings(),
		//	"modulename":     modulename,
		//	"controllername": controllername,
		//	"actionname":     actionname,
		//	"jsname":         fmt.Sprintf("frontend/%s", controllername),
		//	"moduleurl":      fmt.Sprintf("/%s", modulename),
		//	"language":       utils.Langset(c.Request.Header.Get("Accept-Language")),
		//}
		//
		//view.Assign("site", siteViper.AllSettings())
		//view.Assign("config", config)
	}
}
