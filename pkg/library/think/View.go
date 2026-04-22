package think

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type View struct {
}

// Assign 模板变量赋值
// 参数:
//
//	变量名: name 支持 string 或 map[string]any
//	变量值: value  当 name 为 string 时必需
//
// 返回值:
//
//	*View
func (t *View) Assign(c *gin.Context, name any, value ...any) *View {
	sm := c.GetStringMap("__data__")
	if sm == nil {
		sm = map[string]any{
			"_": c,
		}
	}
	switch v := name.(type) {
	case map[string]any:
		for k, val := range v {
			sm[k] = val
		}
	case string:
		var val any
		if len(value) > 0 {
			val = value[0]
		}
		sm[v] = val
	}
	c.Set("__data__", sm)
	return t
}

// Fetch 解析和获取模板内容 用于输出
// 参数:
//
//	模板文件名或者内容: template
//	模板输出变量: vars
//
// 返回值:
//
//	*View
func (t *View) Fetch(c *gin.Context, args ...any) {
	template := ""
	obj := c.GetStringMap("__data__")

	// 解析模板参数
	if len(args) > 0 {
		if templ, ok := args[0].(string); ok {
			template = templ
		}
	}
	// 自动生成模板路径
	if template == "" {
		modulename, _ := c.Get("modulename")
		controllername, _ := c.Get("controllername")
		actionname, _ := c.Get("actionname")
		template = fmt.Sprintf("%s/view/%s/%s.html", modulename, controllername, actionname)
	}
	// 合并额外传入的数据
	if len(args) > 1 {
		if o, ok := args[1].(map[string]any); ok {
			if obj != nil {
				for k, v := range o {
					obj[k] = v // 额外的数据会覆盖或补充 __data__ 中的数据
				}
			} else {
				obj = o
			}
		}
	}
	c.HTML(http.StatusOK, template, obj)
}
