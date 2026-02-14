package multi

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	tTemplate "text/template"
	"time"

	"github.com/gin-gonic/gin/render"
)

var globalFuncMap = template.FuncMap{
	"htmlentities": tTemplate.HTMLEscaper, // 将字符转换为 HTMLProduction 转义字符
	"html_entity_decode": func(str string) template.HTML { // 字符串转化为html标签
		return template.HTML(str)
	},
	"timestamp": func(t ...any) int64 {
		if len(t) > 0 {
			if tm, ok := t[0].(time.Time); ok {
				return tm.Unix()
			}
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
}

// 设置最大递归深度
const maxDepth = 10

// 模板匹配正则表达式，用于提取模板中的子模板引用
var templateRegex = regexp.MustCompile(`{{ *template +["'](.*?)(?:["'].*?)?}}`)

// IRender 模板引擎公共接口
type IRender interface {
	LoadHTMLGlob(string)
	Instance(string, interface{}) render.Render
	SetPrefix(string)
	SetDir(fs.FS)
	SetFuncMap(template.FuncMap)
	SetReplaces(map[string]string)
	SetDelims(render.Delims)
	SetScopeFuncMap(ScopeFuncMap)
	readFileOS(string) (string, string, error)
}

// Tpl 模板结构体，包含模板名称和内容
type Tpl struct {
	Name    string // 模板名称
	Content string // 模板内容
}
type Option func(IRender)

type ScopeFuncMap func(string, template.FuncMap) (string, template.FuncMap)

func WithPrefix(prefix string) Option {
	return func(r IRender) {
		r.SetPrefix(prefix)
	}
}
func WithDir(dir fs.FS) Option {
	return func(r IRender) {
		r.SetDir(dir)
	}
}
func WithDelims(delims render.Delims) Option {
	return func(r IRender) {
		r.SetDelims(delims)
	}
}
func WithReplaces(replaces map[string]string) Option {
	return func(r IRender) {
		r.SetReplaces(replaces)
	}
}

func WithFuncMap(funcMap template.FuncMap) Option {
	return func(r IRender) {
		r.SetFuncMap(funcMap)
	}
}
func WithScopeFuncMap(scopeFuncMap ScopeFuncMap) Option {
	return func(r IRender) {
		r.SetScopeFuncMap(scopeFuncMap)
	}
}

func Render(debug bool, opts ...Option) (r IRender) {
	defDelims := render.Delims{Left: "{{", Right: "}}"}
	defDir := os.DirFS("./internal")
	if debug {
		r = &HTMLDebug{
			delims:  defDelims,
			dir:     defDir,
			funcMap: make(template.FuncMap),
		}
	} else {
		r = &HTMLProduction{
			templates: make(map[string]*template.Template),
			delims:    defDelims,
			dir:       defDir,
			funcMap:   make(template.FuncMap),
		}
	}
	for _, opt := range opts {
		opt(r)
	}
	return
}

// 收集模板
func collect(str string, readFileOS func(string) (string, string, error)) (ts []*Tpl) {
	visited := make(map[string]bool)
	return deep(str, readFileOS, visited, 0, maxDepth)
}

// deepDep 内部递归函数，使用访问记录和深度限制防止循环
func deep(str string, readFileOS func(string) (string, string, error), visited map[string]bool, currentDepth, maxDepth int) (ts []*Tpl) {
	// 检查递归深度
	if currentDepth >= maxDepth {
		return
	}

	// 查找模板中引用的子模板
	templateMatches := templateRegex.FindAllStringSubmatch(str, -1)
	for _, match := range templateMatches {
		if len(match) > 1 {
			filename := match[1] // 获取子模板文件名

			// 检查是否已访问过此模板，防止循环引用
			if visited[filename] {
				continue
			}

			visited[filename] = true // 标记为已访问

			name, s, err := readFileOS(filename)
			if err != nil {
				// 如果文件读取失败，移除访问标记，允许其他路径尝试
				delete(visited, filename)
				continue
			}

			// 添加当前子模板
			ts = append(ts, &Tpl{
				Name:    name,
				Content: s,
			})

			// 递归查找子模板的依赖（增加深度计数）
			childDeps := deep(s, readFileOS, visited, currentDepth+1, maxDepth)
			ts = append(ts, childDeps...)
		}
	}
	return
}

func readFileOS(dir fs.FS, file string, replaces map[string]string) (name string, s string, err error) {
	name = filepath.ToSlash(file) // 统一路径分隔符
	f, err := dir.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return
	}
	s = string(b)
	// 根据替换映射替换内容
	for o, n := range replaces {
		s = strings.ReplaceAll(s, strings.ToUpper(o), n)
	}
	return
}
