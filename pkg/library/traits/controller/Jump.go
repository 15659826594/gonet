package controller

import (
	Config "gota/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ResponseTypeJSON   = "json"
	ResponseTypeXML    = "xml"
	ResponseTypeJSONP  = "jsonp"
	DefaultCodeSuccess = 1
	DefaultCodeError   = 0
)

type Jump struct {
}

type Result struct {
	Code   *int              `json:"code"`
	Msg    string            `json:"msg"`
	Data   any               `json:"data"`
	Url    string            `json:"url"`
	Wait   *int              `json:"wait"`
	Header map[string]string `json:"-"`
}

func (t *Result) Response(c *gin.Context, tmpl string) {
	if t.Header != nil {
		for k, v := range t.Header {
			c.Header(k, v)
		}
	}
	wait := 3
	if t.Wait != nil {
		wait = *t.Wait
	}
	c.HTML(http.StatusOK, tmpl, gin.H{
		"code": *t.Code,
		"msg":  t.Msg,
		"data": t.Data,
		"url":  t.Url,
		"wait": wait,
	})
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
func (j *Jump) Success(c *gin.Context, args ...any) {
	result := defaultResult(args)
	defaultCode := DefaultCodeSuccess
	result.Code = &defaultCode
	result.Response(c, Config.Viper().DispatchSuccessTmpl)
	panic(struct{}{})
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
	result := defaultResult(args...)
	defaultCode := DefaultCodeError
	result.Code = &defaultCode
	result.Response(c, Config.Viper().DispatchErrorTmpl)
	panic(struct{}{})
}

// Result 返回封装后的 API 数据到客户端
// 参数:
//
//	要返回的数据: data
//	返回的 code: code
//	提示信息: msg
//	返回数据格式: type
//	发送的 Header 信息: header
//
// 返回值:
//
//	void
//
// throws HttpResponseException
func (j *Jump) Result(c *gin.Context, data any, code *int, msg string, types string, header map[string]string) {
	panic(struct{}{})
}

func defaultResult(args ...any) *Result {
	var msg string
	var url string
	var data any
	var wait *int
	header := make(map[string]string)
	if len(args) > 0 {
		if arg0, ok := args[0].(string); ok {
			msg = arg0
		}
	}
	if len(args) > 1 {
		if arg1, ok := args[1].(string); ok {
			url = arg1
		}
	}
	if len(args) > 2 {
		data = args[2]
	}
	if len(args) > 3 {
		if arg3, ok := args[3].(*int); ok {
			wait = arg3
		}
	}
	if len(args) > 4 {
		if arg4, ok := args[4].(map[string]string); ok {
			header = arg4
		}
	}
	return &Result{
		Msg:    msg,
		Data:   data,
		Url:    url,
		Wait:   wait,
		Header: header,
	}
}
