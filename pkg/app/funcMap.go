package app

import (
	"fmt"
	"gota/pkg/i18n"
	"gota/pkg/utils"
	"html/template"
	"net/url"
	"strings"
	textTemplate "text/template"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var funcMap = template.FuncMap{
	"json_encode": utils.JsonEncode,
	"json_decode": utils.JsonDecode,
	"urlencode":   url.QueryEscape,
	"cdnurl": func(url string) string {
		return url
	},
	"htmlentities": textTemplate.HTMLEscaper, // 将字符转换为 HTMLProduction 转义字符
	"html_entity_decode": func(str string) template.HTML { // 字符串转化为html标签
		return template.HTML(str)
	},
	"timestamp": func(args ...time.Time) int64 {
		if len(args) > 0 {
			return args[0].Unix()
		}
		return time.Now().Unix()
	},
	"date": func(t int64, format string) string {
		if format == "" {
			format = time.DateTime
		}
		if t == 0 {
			return time.Now().Format(format)
		}
		return time.Unix(t, 0).Format(format)
	},
	"default": func(defaultValue string, value any) string {
		if value == nil || value == "" {
			return defaultValue
		}
		if s, ok := value.(string); ok && s == "" {
			return defaultValue
		}
		return fmt.Sprintf("%v", value)
	},
	"token": func(context *gin.Context, args ...string) template.HTML {
		name := "__token__"
		if len(args) > 0 {
			name = args[0]
		}
		format := "<input type=\"hidden\" name=\"%s\" value=\"%s\" />"
		session := sessions.Default(context)
		val := session.Get(name)
		if s, ok := val.(string); ok {
			return template.HTML(fmt.Sprintf(format, name, s))
		}
		token := utils.Md5(name)
		session.Set(name, token)
		_ = session.Save()
		return template.HTML(fmt.Sprintf(format, name, token))
	},
}

// I18n 子模板 独立作用域
func I18n(name string) any {
	return func(messageID string, templates ...map[string]any) string {
		return i18n.T(name, messageID, templates...)
	}
}

func Url(name string) any {
	return func(url string, args ...any) string {
		var vars string
		if len(args) > 0 {
			switch arg0 := args[0].(type) {
			case string:
				vars = arg0
			case map[string]any:
				parts := make([]string, 0, len(arg0))
				for k, v := range arg0 {
					parts = append(parts, fmt.Sprintf("%s=%v", k, v))
				}
				vars = strings.Join(parts, "&")
			}
		}
		return new(utils.Url).Build(name, url, vars)
	}
}
