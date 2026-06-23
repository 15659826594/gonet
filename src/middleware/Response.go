package middleware

import (
	"github.com/gin-gonic/gin"
)

// ResponseHandler 捕获正常返回
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if _, ok := rec.(struct{}); ok {
				c.Abort()
				return
			}
			panic(rec)
		}()
		c.Next()
	}
}
