package multi

import (
	"html/template"
	"io/fs"
	"iter"

	"github.com/gin-gonic/gin/render"
)

var _ IRender = (*HTMLProduction)(nil)

// HTMLProduction 模板渲染引擎结构体
type HTMLProduction struct {
	html      map[string]*template.Template
	templates []*Tpl           // 已加载的模板集合
	dir       fs.FS            // 文件系统接口，用于读取模板文件
	FuncMap   template.FuncMap // 模板函数映射
}

// LoadHTMLFile 加载HTML模板文件
func (r *HTMLProduction) LoadHTMLFile(filename string) {
	templ := loadHTMLFile(filename)
	if templ != nil {
		r.templates = append(r.templates, templ)
	}
}

// LoadHTMLGlob 加载匹配模式的所有HTML模板文件
func (r *HTMLProduction) LoadHTMLGlob(fs fs.FS, pattern string) {
	templs := loadHTMLGlob(fs, pattern)
	if len(templs) != 0 {
		r.templates = append(r.templates, templs...)
	}
}

func (r *HTMLProduction) Seq2() iter.Seq2[string, *Tpl] {
	return seq2(r.templates)
}

func (r *HTMLProduction) loadTemplate(name string) *template.Template {
	tpl, err := searchSlices(r.templates, name)

	if err != nil {
		panic(err)
	}

	tmpl, err := DOMParser(nil, tpl.Name, r.templates, r.FuncMap)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func (r *HTMLProduction) Init() IRender {
	for _, name := range r.templates {
		r.html[name.Name] = r.loadTemplate(name.Name)
	}
	return r
}

// Instance 创建HTML渲染实例
// name: 模板名称
// data: 渲染数据
func (r *HTMLProduction) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: r.html[name], // 使用预加载的模板
		Data:     data,         // 渲染数据
	}
}
