package multi

import (
	"fmt"
	"gota/pkg/template/parse"
	"html/template"
	"io"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/gin-gonic/gin/render"
)

var (
	globalFuncMap  template.FuncMap
	globalReplaces map[string]string
	globalDelims   = render.Delims{Left: "{{", Right: "}}"}
)

func SetFuncMap(funcMap template.FuncMap) {
	globalFuncMap = funcMap
}

func SetReplaces(replaces map[string]string) {
	globalReplaces = replaces
}

func SetDelims(delims render.Delims) {
	globalDelims = delims
}

func New(debug bool) IRender {
	funcMap := globalFuncMap

	if debug {
		return &HTMLDebug{
			Prefix:  "[GIN-debug] ",
			FuncMap: funcMap,
		}
	}
	return &HTMLProduction{
		html:    make(map[string]*template.Template),
		FuncMap: funcMap,
	}
}

// IRender 模板引擎公共接口
type IRender interface {
	Init() IRender
	LoadHTMLFile(string)
	LoadHTMLGlob(fs.FS, string)
	Seq2() iter.Seq2[string, *Tpl]
	Instance(string, any) render.Render
}

// Tpl 模板结构体，包含模板名称和内容
type Tpl struct {
	dir     fs.FS
	Name    string // 模板名称
	funcMap template.FuncMap
	once    sync.Once
}

func (t *Tpl) SetFuncMap(name string, fun any) {
	t.once.Do(func() {
		if t.funcMap == nil {
			t.funcMap = make(template.FuncMap)
		}
	})
	t.funcMap[name] = fun
}

// DOMParser 解析html文件
func DOMParser(templ *template.Template, file string, templates []*Tpl, funcMap template.FuncMap, chains ...[]string) (*template.Template, error) {
	// 检测循环依赖
	chain, err := depCheck(file, chains)
	if err != nil {
		return nil, err
	}

	tpl, err := searchSlices(templates, file)

	if err != nil {
		return nil, err
	}

	name, s, err := readDirOS(tpl.dir, file, globalReplaces)

	if err != nil {
		return nil, fmt.Errorf("%s：%v\n", err, file)
	}

	if templ == nil {
		templ = template.New(name)
	} else {
		templ = templ.New(name)
	}

	templ.Delims(globalDelims.Left, globalDelims.Right).Funcs(mergeSlices(globalFuncMap, funcMap, tpl.funcMap))

	//分析模板继承关系
	tree, err := parse.Parse(name, s, globalDelims.Left, globalDelims.Right)

	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v\n", err)
	}

	var parseErr error
	for _, t := range tree {
		inspect(t.Root, func(n parse.Node) bool {
			if n == nil {
				return false
			}
			switch node := n.(type) {
			case *parse.BlockNode:
				break
			case *parse.TemplateNode:
				_, e := DOMParser(templ, node.Name, templates, funcMap, chain)
				if e != nil {
					parseErr = e
					return false
				}
			}
			return true
		})

		if parseErr != nil {
			break
		}
	}

	if parseErr != nil {
		return nil, parseErr
	}

	_, err = templ.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v\n", err)
	}

	return templ, nil
}

// depCheck 检测模板文件是否存在循环依赖
func depCheck(file string, chains [][]string) ([]string, error) {
	var chain []string
	if len(chains) > 0 {
		chain = make([]string, len(chains[0]))
		copy(chain, chains[0])
	}

	for _, calledFile := range chain {
		if calledFile == file {
			cyclePath := strings.Join(append(chain, file), " -> ")
			return nil, fmt.Errorf("circular template dependency detected: %s", cyclePath)
		}
	}

	chain = append(chain, file)
	return chain, nil
}

// inspect 获取模板节点
func inspect(node parse.Node, f func(parse.Node) bool) {
	if !f(node) {
		return
	}
	switch n := node.(type) {
	case *parse.ListNode:
		if n.Nodes != nil {
			for _, child := range n.Nodes {
				inspect(child, f)
			}
		}
	}
}

// 加载文件
func loadHTMLFile(filename string) *Tpl {
	tpl := &Tpl{
		dir:  nil,
		Name: filename,
	}
	return tpl
}

// 加载文件夹下的文件 , 支持模糊匹配
func loadHTMLGlob(fs fs.FS, pattern string) []*Tpl {
	filenames, err := doublestar.Glob(fs, pattern)
	if err != nil {
		panic(err)
	}
	if len(filenames) == 0 {
		panic(fmt.Errorf("html/template: pattern matches no files: %#q", pattern))
	}
	var templates []*Tpl
	for _, filename := range filenames {
		templates = append(templates, &Tpl{
			dir:  fs,
			Name: filename,
		})
	}
	return templates
}

// 遍历已加载的模板
func seq2(templates []*Tpl) iter.Seq2[string, *Tpl] {
	return func(yield func(string, *Tpl) bool) {
		for _, tpl := range templates {
			if !yield(tpl.Name, tpl) {
				return
			}
		}
	}
}

// 读取模板文件
func readDirOS(dir fs.FS, file string, replacess ...map[string]string) (name string, s string, err error) {
	name = filepath.ToSlash(file) // 统一路径分隔符
	var f fs.File
	if dir == nil {
		f, err = os.Open(file)
	} else {
		f, err = dir.Open(file)
	}
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
	if len(replacess) > 0 {
		replaces := replacess[0]
		for o, n := range replaces {
			s = strings.ReplaceAll(s, o, n)
		}
	}
	return
}

func mergeSlices(maps ...template.FuncMap) template.FuncMap {
	merged := make(template.FuncMap)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func searchSlices(templates []*Tpl, name string) (*Tpl, error) {
	for _, t := range templates {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, fmt.Errorf("template not found: %s", name)
}
