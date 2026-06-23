//go:build release

package src

import "github.com/gin-gonic/gin"

func init() {
	EnvGinMode = gin.ReleaseMode
}
