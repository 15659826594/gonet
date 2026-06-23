package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Backend struct {
}

func (t *Backend) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(" Backend Initialize")
	}
}
