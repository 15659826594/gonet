package middleware

import (
	"gonet/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseHandler 捕获正常放回
// . "gonet/internal/common"
// Response(Error("成功消息"))
// Response(Success("失败消息"))
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case *common.Result:
					switch e.Type {
					case "json":
						c.JSON(e.Statuscode, e)
					case "xml":
						c.XML(e.Statuscode, e)
					case "jsonp":
						c.JSONP(e.Statuscode, e)
					default:
						c.JSON(e.Statuscode, e)
					}
				case error:
					c.JSON(http.StatusInternalServerError, common.Error(e.Error()))
				default:
					panic(err)
				}
			}
		}()

		c.Next()
	}
}
