package controller

import "github.com/gin-gonic/gin"

type Api struct {
}

func (c Api) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
