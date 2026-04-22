package multi

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"iter"
	"os"
	"strings"

	"github.com/gin-gonic/gin/render"
)

var _ IRender = (*HTMLDebug)(nil)
var DefaultWriter io.Writer = os.Stdout
var DebugPrintFunc func(format string, values ...any)

// HTMLDebug contains template delims and pattern and function with file list.
type HTMLDebug struct {
	Prefix    string
	templates []*Tpl
	dir       fs.FS
	FuncMap   template.FuncMap
}

func (r *HTMLDebug) LoadHTMLFile(filename string) {
	templ := loadHTMLFile(filename)
	if templ != nil {
		r.templates = append(r.templates, templ)
	}
}

func (r *HTMLDebug) LoadHTMLGlob(fs fs.FS, pattern string) {
	templs := loadHTMLGlob(fs, pattern)
	if len(templs) != 0 {
		r.templates = append(r.templates, templs...)
	}
}

func (r *HTMLDebug) Seq2() iter.Seq2[string, *Tpl] {
	return seq2(r.templates)
}

func (r *HTMLDebug) loadTemplate(name string) *template.Template {
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

func (r *HTMLDebug) debugPrintLoadTemplate() {
	var buf strings.Builder
	for _, templ := range r.templates {
		buf.WriteString("\t- ")
		buf.WriteString(templ.Name)
		buf.WriteString("\n")
	}
	r.debugPrint("Loaded HTMLProduction Templates (%d): \n%s\n", len(r.templates), buf.String())
}

func (r *HTMLDebug) debugPrint(format string, values ...any) {
	if DebugPrintFunc != nil {
		DebugPrintFunc(format, values...)
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(DefaultWriter, r.Prefix+format, values...)
}

func (r *HTMLDebug) Init() IRender {
	r.debugPrintLoadTemplate()
	return r
}

func (r *HTMLDebug) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: r.loadTemplate(name),
		Data:     data,
	}
}
