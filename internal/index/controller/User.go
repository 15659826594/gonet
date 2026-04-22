package controller

import (
	"fmt"
	"gota/internal/common/controller"
	"gota/internal/common/library/Auth"
	"gota/pkg/app/route"

	"github.com/gin-contrib/sessions"
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

// Index 会员中心
func (t User) Index(c *gin.Context) {
	t.View.Assign(c, "title", "User center")
	t.View.Fetch(c)
}

// Register 注册会员
func (t User) Register(c *gin.Context) {
	url := c.DefaultQuery("url", "我是登录")
	t.Assign(c, "url", url)
	t.Assign(c, "title", t.T(c, "Login"))
	t.View.Fetch(c)
}

type LoginForm struct {
	Account   string `form:"account" binding:"required,min=3,max=50"`
	Password  string `form:"password" binding:"required,min=6,max=30"`
	KeepLogin int    `form:"keeplogin"`
	Token     string `form:"__token__" binding:"required"`
}

// Login 会员登录
func (t User) Login(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println(session.Get("name"))
	url := c.DefaultQuery("url", "我是登录")
	if c.Request.Method == "POST" {
		var form LoginForm
		if err := c.ShouldBind(&form); err != nil {
			t.Error(c, t.T(c, err.Error()), nil, gin.H{"token": t.Frontend.Token(c)})
		}

		auth := Auth.Instance(c)
		if auth.Login(form.Account, form.Password) {
			t.Success(c, t.T(c, "Logged in successful"))
		} else {
			t.Error(c, t.T(c, auth.GetError().Error()), nil, gin.H{"token": t.Frontend.Token(c)})
		}
	}
	t.Assign(c, "url", url)
	t.Assign(c, "title", t.T(c, "Login"))
	t.View.Fetch(c)
}

// Logout 退出登录
func (t User) Logout(c *gin.Context) {
}

// Profile 个人信息
func (t User) Profile(c *gin.Context) {
	t.Assign(c, "title", t.T(c, "Profile"))
	t.View.Fetch(c)
}

// Changepwd 修改密码
func (t User) Changepwd(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Assign("title", i18n.T(c.GetString("url"), "Change password"))
	//tpl.Display("index/view/user/login.html")
}

func (t User) Attachment(c *gin.Context) {
	//tpl := c.MustGet("Think").(*template.Template)
	//tpl.Display("index/view/user/attachment.html")
}
