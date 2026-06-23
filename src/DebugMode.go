//go:build debug

package src

import "github.com/gin-gonic/gin"

func init() {
	EnvGinMode = gin.DebugMode
}
