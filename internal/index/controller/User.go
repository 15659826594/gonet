package controller

import (
	"gonet/internal/common/controller"
	"gonet/pkg/app/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	route.Register(User{
		NoNeedLogin: []string{"login", "register", "third"},
		NoNeedRight: []string{"*"},
	})
}

type User struct {
	controller.Frontend
	NoNeedLogin []string
	NoNeedRight []string
}

func (u User) Index(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Assign("title", i18n.T(c.GetString("url"), "User center"))
	//tpl.Display("index/view/user/index.html")
	c.String(http.StatusOK, "11")
}

func (u User) Register(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Display("index/view/user/register.html")
}

func (u User) Login(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Assign("title", i18n.T(c.GetString("url"), "Login"))
	//tpl.Display("index/view/user/login.html")
}

func (u User) Logout(c *gin.Context) {
}

func (u User) Profile(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Assign("title", i18n.T(c.GetString("url"), "Profile"))
	//tpl.Display("index/view/user/login.html")
}

func (u User) Changepwd(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Assign("title", i18n.T(c.GetString("url"), "Change password"))
	//tpl.Display("index/view/user/login.html")
}

func (u User) Attachment(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Display("index/view/user/attachment.html")
}
