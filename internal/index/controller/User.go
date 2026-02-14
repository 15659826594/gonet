package controller

import (
	"fmt"
	"gonet/pkg/app/route"
	Config "gonet/pkg/config"
	"gonet/pkg/i18n"
	"gonet/pkg/template"
	"gonet/pkg/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	route.Register(User{
		NoNeedLogin: []string{"login", "register", "third"},
		NoNeedRight: []string{"*"},
	})
}

type User struct {
	NoNeedLogin []string
	NoNeedRight []string
}

func (u User) BeforeAction() []gin.HandlerFunc {
	return []gin.HandlerFunc{func(c *gin.Context) {
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
	}}
}

func (u User) Index(c *gin.Context) {
	tpl := c.MustGet("Think").(*template.Template)
	tpl.Assign("title", i18n.T(c.GetString("url"), "User center"))
	tpl.Display("index/view/user/index.html")
}

func (u User) Register(c *gin.Context) {
	tpl := c.MustGet("Think").(*template.Template)
	tpl.Display("index/view/user/register.html")
}

func (u User) Login(c *gin.Context) {
	tpl := c.MustGet("Think").(*template.Template)
	tpl.Assign("title", i18n.T(c.GetString("url"), "Login"))
	tpl.Display("index/view/user/login.html")
}

func (u User) Logout(c *gin.Context) {
}

func (u User) Profile(c *gin.Context) {
	tpl := c.MustGet("Think").(*template.Template)
	tpl.Assign("title", i18n.T(c.GetString("url"), "Profile"))
	tpl.Display("index/view/user/login.html")
}

func (u User) Changepwd(c *gin.Context) {
	tpl := c.MustGet("Think").(*template.Template)
	tpl.Assign("title", i18n.T(c.GetString("url"), "Change password"))
	tpl.Display("index/view/user/login.html")
}

func (u User) Attachment(c *gin.Context) {
	tpl := c.MustGet("Think").(*template.Template)
	tpl.Display("index/view/user/attachment.html")
}
