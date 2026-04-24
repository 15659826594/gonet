package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	store := cookie.NewStore([]byte("your-secret-key-must-be-at-least-32-bytes-long!!"))
	store.Options(sessions.Options{
		Path:     "/",                  // Cookie 路径
		MaxAge:   86400,                // Cookie 有效期（秒），这里是 24 小时
		HttpOnly: true,                 // 仅 HTTP 可访问，防止 XSS 攻击
		Secure:   false,                // 生产环境如果使用 HTTPS 应设为 true
		SameSite: http.SameSiteLaxMode, // 防止 CSRF 攻击
	})
	return sessions.Sessions("session", store)
}
