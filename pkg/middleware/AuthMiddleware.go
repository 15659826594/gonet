package middleware

import (
	. "gonet/internal/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

}

func FastAuthMiddleware(c *gin.Context) {
	token := GinContextGetToken(c)
	if token == "" {
		Response(Error("请登陆后操作"))
	}
	user, err := GetUserByToken(token)
	if err != nil {
		Response(Error(err.Error()))
	}
	c.Set("auth", user)
	c.Next()
}

func GinContextGetToken(c *gin.Context) string {
	var token string
	// 1. 从Header中获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 支持 Bearer token 格式
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = authHeader
		}
	}
	// 2. 如果Header中没有token，从GET参数中获取
	if token == "" {
		token = c.Query("token")
	}
	// 3. 如果GET参数中没有token，从POST表单中获取
	if token == "" {
		token = c.PostForm("token")
	}
	return token
}
