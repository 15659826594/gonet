package middleware

import (
	"errors"
	"gota/pkg/exception"

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
			if err, ok := rec.(error); ok {
				if ok := errors.Is(err, exception.HttpResponseException); ok {
					c.Abort()
					return
				}
			}
			panic(rec)
		}()
		c.Next()
	}
}
