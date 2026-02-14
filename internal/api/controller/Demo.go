package controller

import (
	. "gonet/internal/common"

	"github.com/gin-gonic/gin"
)

/*Demo
 * 示例接口
 */
type Demo struct {
	//如果$noNeedLogin为空表示所有接口都需要登录才能请求
	//如果$noNeedRight为空表示所有接口都需要验证权限才能请求
	//如果接口已经设置无需登录,那也就无需鉴权了
	//
	// 无需登录的接口,*表示全部
	NoNeedLogin []string
	// 无需鉴权的接口,*表示全部
	NoNeedRight []string
}

// Test
// @Summary 测试方法
// @Description 测试描述信息
// @Tags demo
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} Result "发送成功"
// @Router /demo/test [post]
// @Security ApiKeyAuth
func (Demo) Test(c *gin.Context) {
	Response(Success("返回成功"))
}

// Test1
// @Summary 无需登录的接口
// @Description 无需登录的接口
// @Tags demo
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} Result "发送成功"
// @Router /demo/test1 [post]
// @Security ApiKeyAuth
func (Demo) Test1(c *gin.Context) {
	Response(Success("返回成功", gin.H{
		"action": "test1",
	}))
}

// Test2
// @Summary 需要登录的接口
// @Description 需要登录的接口
// @Tags demo
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} Result "发送成功"
// @Router /demo/test2 [post]
// @Security ApiKeyAuth
func (Demo) Test2(c *gin.Context) {
	Response(Success("返回成功", gin.H{
		"action": "test2",
	}))
}

// Test3
// @Summary 需要登录且需要验证有相应组的权限
// @Description 需要登录且需要验证有相应组的权限
// @Tags demo
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} Result "发送成功"
// @Router /demo/test3 [get]
// @Security ApiKeyAuth
func (Demo) Test3(c *gin.Context) {
	Response(Success("返回成功", gin.H{
		"action": "test3",
	}))
}
