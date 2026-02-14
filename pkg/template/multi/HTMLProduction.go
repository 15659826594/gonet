package multi

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/gin-gonic/gin/render"
)

// HTMLProduction 模板渲染引擎结构体
type HTMLProduction struct {
	templates    map[string]*template.Template // 已加载的模板集合
	dir          fs.FS                         // 文件系统接口，用于读取模板文件
	replaces     map[string]string             // 模板内容替换映射
	delims       render.Delims                 // 模板分隔符配置（如 {{ 和 }}）
	funcMap      template.FuncMap              // 模板函数映射
	scopeFuncMap ScopeFuncMap
}

func (r *HTMLProduction) SetPrefix(string) {}

func (r *HTMLProduction) SetDir(dir fs.FS) {
	r.dir = dir
}
func (r *HTMLProduction) SetReplaces(replaces map[string]string) {
	r.replaces = replaces
}
func (r *HTMLProduction) SetDelims(delims render.Delims) {
	r.delims = delims
}
func (r *HTMLProduction) SetFuncMap(funcMap template.FuncMap) {
	r.funcMap = funcMap
}

func (r *HTMLProduction) SetScopeFuncMap(scopeFuncMap ScopeFuncMap) {
	r.scopeFuncMap = scopeFuncMap
}

// LoadHTMLGlob 加载匹配模式的所有HTML模板文件
func (r *HTMLProduction) LoadHTMLGlob(pattern string) {
	_, err := r.parseGlob(pattern)
	if err != nil {
		panic(err)
	}
}

// LoadHTML 加载HTML模板文件
func (r *HTMLProduction) LoadHTML(filename string) {
	_, err := r.parseFiles(r.readFileOS, filename)
	if err != nil {
		panic(err)
	}
	var buf strings.Builder
	buf.WriteString("\t- ")
	buf.WriteString(filename)
	buf.WriteString("\n")
	fmt.Println(buf.String())
}

// LoadTemplate 创建新的模板实例
// name: 模板名称
// delims: 模板分隔符
// funcMap: 模板函数映射
// tpls: 模板列表
func (r *HTMLProduction) LoadTemplate(html *Tpl, delims render.Delims, funcMap template.FuncMap, tpls ...*Tpl) (*template.Template, error) {
	// 创建新模板并设置分隔符和函数映射
	name := html.Name
	scope := make(template.FuncMap)
	if r.scopeFuncMap != nil {
		name, scope = r.scopeFuncMap(name, scope)
	}
	tmpl := template.New(name).Delims(delims.Left, delims.Right).Funcs(globalFuncMap).Funcs(funcMap).Funcs(scope)
	for _, tpl := range tpls {
		// 为每个子模板创建并解析
		_, err := tmpl.New(tpl.Name).Parse(tpl.Content)
		if err != nil {
			return nil, err
		}
	}
	_, err := tmpl.Parse(html.Content)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// parseFiles 解析多个模板文件
// readFile: 读取文件的函数
// filenames: 文件名列表
func (r *HTMLProduction) parseFiles(readFile func(string) (string, string, error), filenames ...string) (map[string]*template.Template, error) {
	if len(filenames) == 0 {
		// 没有指定文件时返回错误
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		name, s, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		// 解析当前模板及其依赖的子模板
		tmpl, err := r.LoadTemplate(&Tpl{
			Name:    name,
			Content: s,
		}, r.delims, r.funcMap, collect(s, r.readFileOS)...)
		if err != nil {
			return nil, err
		}
		// 将解析后的模板存入引擎
		r.templates[name] = tmpl
	}
	return r.templates, nil
}

// parseGlob 解析匹配模式的所有文件
func (r *HTMLProduction) parseGlob(pattern string) (map[string]*template.Template, error) {
	// 使用 doublestar 库进行文件名模式匹配
	filenames, err := doublestar.Glob(r.dir, pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: pattern matches no files: %#q", pattern)
	}
	// 解析找到的所有文件
	return r.parseFiles(r.readFileOS, filenames...)
}

// 模板缓存，避免重复读取文件
var tempCache = make(map[string]string)

// readFileOS 从文件系统读取文件内容
func (r *HTMLProduction) readFileOS(file string) (name string, s string, err error) {
	name = filepath.ToSlash(file)
	if t, exist := tempCache[name]; exist {
		return name, t, nil
	}
	name, s, err = readFileOS(r.dir, file, r.replaces)
	tempCache[name] = s
	return
}

// Instance 创建HTML渲染实例
// name: 模板名称
// data: 渲染数据
func (r *HTMLProduction) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.templates[name], // 使用预加载的模板
		Data:     data,              // 渲染数据
	}
}
