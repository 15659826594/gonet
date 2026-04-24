package controller

import (
	"fmt"
	"gota/internal/common/library/Auth"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const (
	ResponseTypeJSON   = "json"
	ResponseTypeXML    = "xml"
	ResponseTypeJSONP  = "jsonp"
	DefaultCodeSuccess = 1
	DefaultCodeError   = 0
)

type Api struct {
}

func (t *Api) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Auth.Instance(c)

		modulename := c.GetString("modulename")
		controllername := c.GetString("controllername")
		actionname := c.GetString("actionname")

		//parseToken
		token := parseToken(c)

		path := strings.ReplaceAll(controllername, ".", "/") + "/" + actionname

		// 设置当前请求的URI
		auth.SetRequestUri(path)
		// 检测是否需要验证登录
		if !auth.Match(c.GetStringSlice("noNeedLogin")) {
			//初始化
			auth.Init(token)
			//检测是否登录
			if auth.IsLogin() {
				t.Error(c, "Please login first", "index/user/login")
			}
			// 判断是否需要验证权限
			if !auth.Match(c.GetStringSlice("noNeedRight")) {
				// 判断控制器和方法判断是否有对应权限
				if auth.Check(path, modulename) {
					t.Error(c, "You have no permission")
				}
			}
		} else {
			// 如果有传递token才验证是否登录状态
			if token != "" {
				auth.Init(token)
			}
		}
		fmt.Println(fmt.Sprintf("token是%s", token))
	}
}

// Success 操作成功返回的数据
// 参数:
//
//	提示信息: msg
//	要返回的数据: data
//	错误码，默认为1: code
//	输出类型: type
//	发送的 Header 信息: header
func (t *Api) Success(c *gin.Context, args ...any) {
	msg, data, code, types, header := defaultResult(args...)
	if code == nil {
		defaultCode := DefaultCodeSuccess
		code = &defaultCode
	}
	t.Result(c, msg, data, code, types, header)
}

// Error 操作失败返回的数据
// 参数:
//
//	提示信息: msg
//	要返回的数据: data
//	错误码，默认为0: code
//	输出类型: type
//	发送的 Header 信息: header
func (t *Api) Error(c *gin.Context, args ...any) {
	msg, data, code, types, header := defaultResult(args...)
	if code == nil {
		defaultCode := DefaultCodeError
		code = &defaultCode
	}
	t.Result(c, msg, data, code, types, header)
}

func defaultResult(args ...any) (string, any, *int, string, map[string]string) {
	var msg string
	var data any
	var code *int
	var types string
	header := make(map[string]string)
	if len(args) > 0 {
		msg, _ = args[0].(string)
	}
	if len(args) > 1 {
		data = args[1]
	}
	if len(args) > 2 {
		if c, ok := args[2].(*int); ok {
			code = c
		}
	}
	if len(args) > 3 {
		types, _ = args[3].(string)
	}
	if len(args) > 4 {
		if h, ok := args[4].(map[string]string); ok {
			header = h
		}
	}
	return msg, data, code, types, header
}

// Result 返回封装后的 API 数据到客户端
// 参数:
//
//	提示信息: msg
//	要返回的数据: data
//	错误码，默认为0: code
//	输出类型，支持json/xml/jsonp: type
//	发送的 Header 信息: header
//
// 返回值:
//
//	void
//
// throws HttpResponseException
func (*Api) Result(c *gin.Context, msg string, data any, code *int, t string, header map[string]string) {
	codeValue := 0
	if code != nil {
		codeValue = *code
	}
	result := map[string]any{
		"code": codeValue,
		"msg":  msg,
		"time": time.Now().Unix(),
		"data": data,
	}
	if t == "" {
		t = ResponseTypeJSON
	}
	ResponseCreate(c, result, t, code, header)
	panic(struct{}{})
}

func ResponseCreate(c *gin.Context, result map[string]any, t string, code *int, header map[string]string) {
	// 获取或设置 HTTP 状态码
	statusCode := 200
	if code != nil {
		// 根据 code 值判断 HTTP 状态码
		if *code >= 1000 || *code < 200 {
			statusCode = 200
		} else {
			statusCode = *code
		}
	}

	// 先设置自定义 header（必须在写入响应之前）
	for key, value := range header {
		c.Header(key, value)
	}

	// 设置响应时间
	if gin.IsDebugging() {
		result["latency"] = time.Now().UnixMilli() - c.GetInt64("startTime")
	}

	// 根据不同的响应类型返回数据
	switch t {
	case ResponseTypeJSON:
		c.JSON(statusCode, result)
	case ResponseTypeXML:
		c.XML(statusCode, result)
	case ResponseTypeJSONP:
		callback := c.Query("callback")
		if callback == "" {
			callback = "callback"
		}
		c.Writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		c.JSONP(statusCode, gin.H{
			"callback": callback,
			"data":     result,
		})
	default:
		c.JSON(statusCode, result)
	}
}

// ... existing code ...

// parseToken 从请求中获取 parseToken
// 优先级顺序：
// 1. Authorization Header (Bearer 格式)
// 2. HTTP_TOKEN Header
// 3. URL 查询参数 parseToken
// 4. POST 表单参数 parseToken
// 5. Cookie 中的 parseToken
//
// 参数:
//   - c: gin.Context 上下文对象
//
// 返回值:
//   - string: parseToken 值，如果未找到则返回空字符串
func parseToken(c *gin.Context) string {
	var token string

	// 1. 从 Authorization Header 获取 Bearer Token
	bearerToken := c.GetHeader("Authorization")
	if bearerToken != "" {
		if strings.HasPrefix(bearerToken, "Bearer ") {
			return strings.TrimPrefix(bearerToken, "Bearer ")
		}
		return bearerToken
	}

	// 2. 从 HTTP_TOKEN Header 获取
	token = c.GetHeader("HTTP_TOKEN")
	if token != "" {
		return token
	}

	// 3. 从 URL 查询参数获取
	token = c.Query("parseToken")
	if token != "" {
		return token
	}

	// 4. 从 POST 表单获取
	token = c.PostForm("parseToken")
	if token != "" {
		return token
	}

	// 5. 从 Cookie 获取
	token, _ = c.Cookie("parseToken")
	return token
}

func langset(c *gin.Context) language.Tag {
	tags, _, _ := language.ParseAcceptLanguage(c.Request.Header.Get("Accept-Language"))
	matcher := language.NewMatcher(tags)
	matched, _, _ := matcher.Match()
	return matched
}
