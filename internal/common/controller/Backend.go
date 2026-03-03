package controller

import "github.com/gin-gonic/gin"

type Backend struct {
}

func (c Backend) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
