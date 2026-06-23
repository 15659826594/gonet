//go:build test

package src

import "github.com/gin-gonic/gin"

func init() {
	EnvGinMode = gin.TestMode
}
