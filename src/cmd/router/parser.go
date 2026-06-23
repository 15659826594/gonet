package router

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var (
	// @Router /path [method]
	routerPattern = regexp.MustCompile(`@Router\s+(\S+)\s*\[(\w+)]`)
)

type Ast struct {
	Path        string
	Package     string
	Imports     map[string]*ImportSpec
	StructTypes []*StructType
	FuncDecl    []*FuncDecl
}

type ImportSpec struct {
	Name string
	Path string
	Spec *ast.ImportSpec
}

type StructType struct {
	Name      string
	Docs      string
	FuncDecls []*FuncDecl
}

func NewStructType(nodeType *ast.TypeSpec) (*StructType, error) {
	if _, isStruct := nodeType.Type.(*ast.StructType); !isStruct {
		return nil, fmt.Errorf("type %s is not a struct", nodeType.Name.Name)
	}
	return &StructType{
		Name:      nodeType.Name.Name,
		FuncDecls: make([]*FuncDecl, 0),
	}, nil
}

type FuncDecl struct {
	Name     string
	Docs     string
	Receiver string
	Params   map[string]*ast.Field
	Results  map[int]*ast.Field
}

func NewFuncDecl(nodeType *ast.FuncDecl) (*FuncDecl, error) {
	if nodeType.Name.Name == "init" {
		return nil, fmt.Errorf("skip init function")
	}
	method := &FuncDecl{
		Name: nodeType.Name.Name,
		Docs: nodeType.Doc.Text(),
	}

	if nodeType.Recv != nil && len(nodeType.Recv.List) > 0 {
		recv := nodeType.Recv.List[0]
		switch t := recv.Type.(type) {
		case *ast.Ident:
			method.Receiver = t.Name
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				method.Receiver = "*" + ident.Name
			}
		}
	}

	if nodeType.Type.Params != nil {
		method.Params = make(map[string]*ast.Field, len(nodeType.Type.Params.List))
		for _, field := range nodeType.Type.Params.List {
			for _, name := range field.Names {
				method.Params[name.Name] = field
			}
		}
	}

	if nodeType.Type.Results != nil {
		method.Results = make(map[int]*ast.Field, len(nodeType.Type.Results.List))
		for i, field := range nodeType.Type.Results.List {
			method.Results[i] = field
		}
	}
	return method, nil
}

// NewFileAst 解析单个 Go 文件
func NewFileAst(path string) *Ast {
	return &Ast{
		Path:    path,
		Imports: make(map[string]*ImportSpec),
	}
}

func (p *Ast) ParseFile(filename string) error {
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if p.Package == "" {
		p.Package = file.Name.Name
	}

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		var name string
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			name = parts[len(parts)-1]
		}

		p.Imports[name] = &ImportSpec{
			Name: name,
			Path: path,
			Spec: imp,
		}
	}

	var methods []*FuncDecl
	structMap := make(map[string]string)

	ast.Inspect(file, func(n ast.Node) bool {
		switch nodeType := n.(type) {
		case *ast.GenDecl:
			if nodeType.Tok == token.TYPE {
				for _, spec := range nodeType.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
							obj, _ := NewStructType(typeSpec)
							comment := typeSpec.Doc.Text()
							if comment == "" {
								comment = nodeType.Doc.Text()
							}
							structMap[obj.Name] = comment
						}
					}
				}
			}
		case *ast.TypeSpec:
			if obj, err := NewStructType(nodeType); err == nil {
				p.StructTypes = append(p.StructTypes, obj)
			}
		case *ast.FuncDecl:
			if method, err := NewFuncDecl(nodeType); err == nil {
				methods = append(methods, method)
			}
		}
		return true
	})

	for _, s := range p.StructTypes {
		if docs, ok := structMap[s.Name]; ok {
			s.Docs = docs
		}
		for i := 0; i < len(methods); {
			receiverName := methods[i].Receiver
			// 跳过没有接收器的普通函数
			if receiverName == "" {
				i++
				continue
			}
			// 去掉指针前缀
			if receiverName[0] == '*' {
				receiverName = receiverName[1:]
			}
			if receiverName == s.Name {
				s.FuncDecls = append(s.FuncDecls, methods[i])
				methods = append(methods[:i], methods[i+1:]...)
			} else {
				i++
			}
		}
	}

	p.FuncDecl = append(p.FuncDecl, methods...)
	return nil
}
