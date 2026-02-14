package multi

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/gin-gonic/gin/render"
)

var DefaultWriter io.Writer = os.Stdout
var DebugPrintFunc func(format string, values ...any)

// HTMLDebug contains template delims and pattern and function with file list.
type HTMLDebug struct {
	prefix       string
	Files        []string
	dir          fs.FS
	replaces     map[string]string
	delims       render.Delims
	funcMap      template.FuncMap
	scopeFuncMap ScopeFuncMap
}

func (r *HTMLDebug) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *HTMLDebug) SetDir(dir fs.FS) {
	r.dir = dir
}
func (r *HTMLDebug) SetReplaces(replaces map[string]string) {
	r.replaces = replaces
}
func (r *HTMLDebug) SetDelims(delims render.Delims) {
	r.delims = delims
}
func (r *HTMLDebug) SetFuncMap(funcMap template.FuncMap) {
	r.funcMap = funcMap
}

func (r *HTMLDebug) SetScopeFuncMap(scopeFuncMap ScopeFuncMap) {
	r.scopeFuncMap = scopeFuncMap
}

func (r *HTMLDebug) LoadHTMLGlob(pattern string) {
	filenames, err := doublestar.Glob(r.dir, pattern)
	if err != nil {
		panic(err)
	}
	if len(filenames) == 0 {
		panic(fmt.Sprintf("html/template: pattern matches no files: %#q", pattern))
	}
	r.Files = filenames
	r.debugPrintLoadTemplate()
}

func (r *HTMLDebug) loadTemplate(name string) *template.Template {
	name, s, err := r.readFileOS(name)
	if err != nil {
		return nil
	}
	scope := make(template.FuncMap)
	if r.scopeFuncMap != nil {
		name, scope = r.scopeFuncMap(name, scope)
	}
	tpls := collect(s, r.readFileOS)
	// Funcs 全局方法 | 模板方法 | 子模板独立方法
	tmpl := template.New(name).Delims(r.delims.Left, r.delims.Right).Funcs(globalFuncMap).Funcs(r.funcMap).Funcs(scope)
	for _, tpl := range tpls {
		tmpl.New(tpl.Name).Parse(tpl.Content)
	}
	_, err = tmpl.Parse(s)
	if err == nil {
		return tmpl
	}
	return nil
}

func (r *HTMLDebug) readFileOS(file string) (name string, s string, err error) {
	return readFileOS(r.dir, file, r.replaces)
}

func (r *HTMLDebug) debugPrintLoadTemplate() {
	var buf strings.Builder
	for _, name := range r.Files {
		buf.WriteString("\t- ")
		buf.WriteString(name)
		buf.WriteString("\n")
	}
	r.debugPrint("Loaded HTMLProduction Templates (%d): \n%s\n", len(r.Files), buf.String())
}

func (r *HTMLDebug) debugPrint(format string, values ...any) {
	if DebugPrintFunc != nil {
		DebugPrintFunc(format, values...)
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(DefaultWriter, r.prefix+format, values...)
}

func (r *HTMLDebug) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.loadTemplate(name),
		Data:     data,
	}
}
