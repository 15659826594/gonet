package controller

import (
	"gonet/pkg/exception"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Jump struct {
}

type Result struct {
	Code       int            `json:"code"`
	Msg        string         `json:"msg"`
	Time       int64          `json:"time"`
	CostTime   int64          `json:"costTime,omitempty"`
	Data       any            `json:"data"`
	Statuscode int            `json:"-"`
	Type       string         `json:"-"`
	Header     map[string]any `json:"-"`
}

func (r *Result) Set(args ...any) *Result {
	return nil
}

// Success 操作成功跳转的快捷方法
// 参数:
//
//	提示信息: msg
//	跳转的 URL 地址: url
//	返回的数据: data
//	跳转等待时间: wait
//	发送的 Header 信息: header
//
// 返回值:
//
//	void
//
// throws HttpResponseException
func (j *Jump) Success(c *gin.Context, args ...any) *Result {
	//res := new(Result).Set(args)
	res := &Result{
		Code:       1,
		Msg:        "",
		Time:       time.Now().Unix(),
		Data:       nil,
		Statuscode: http.StatusOK,
		Type:       "json",
	}
	if gin.IsDebugging() {
		res.CostTime = time.Now().Unix() - c.GetInt64("startTime")
	}
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok {
			res.Msg = msg
		}
	}
	if len(args) > 1 {
		res.Data = args[1]
	}
	panic(res)
}

// Error 操作错误跳转的快捷方法
// 参数:
//
//	提示信息: msg
//	跳转的 URL 地址: url
//	返回的数据: data
//	跳转等待时间: wait
//	发送的 Header 信息: header
//
// 返回值:
//
//	void
//
// throws HttpResponseException
func (j *Jump) Error(c *gin.Context, args ...any) {
	res := &Result{
		Code:       0,
		Msg:        "",
		Time:       time.Now().Unix(),
		Data:       nil,
		Statuscode: http.StatusOK,
		Type:       "json",
	}
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok {
			res.Msg = msg
		}
	}
	if len(args) > 1 {
		res.Data = args[1]
	}
	if gin.IsDebugging() {
		res.CostTime = time.Now().UnixMilli() - c.GetInt64("startTime")
	}
	c.JSON(http.StatusOK, res)
	panic(exception.HttpResponseException)
}
