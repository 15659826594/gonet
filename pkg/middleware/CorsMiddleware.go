package middleware

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CorsMiddleware 设置允许跨域
func CorsMiddleware() gin.HandlerFunc {
	// 开发环境全部放行
	if gin.IsDebugging() {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-Token")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Max-Age", "86400") // 预检请求缓存时间（24小时）
			c.Header("Access-Control-Allow-Credentials", "true")
		}
	}
	cors := viper.GetStringSlice("fastadmin.cors_request_domain")
	if len(cors) == 0 {
		cors = []string{"*"}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// 如果没有 Origin 头，说明不是跨域请求，直接放行
		if origin == "" {
			c.Next()
			return
		}

		if !(len(cors) == 1 && cors[0] == "*" || slices.Contains(cors, origin)) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "origin not allowed",
				"data": nil,
				"time": time.Now().Unix(),
			})
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Max-Age", "86400") // 预检请求缓存时间（24小时）

		// 仅当不是通配符 * 时才允许携带 Cookie 或 Authorization 凭据
		if origin != "*" {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
