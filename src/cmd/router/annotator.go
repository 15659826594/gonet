package router

import "text/template"

type AnnotationHandler interface {
	Name() string
	Scope() string
	Parse(comment string, meta *NodeMeta) bool
	Generate(tpl *template.Template, data map[string]interface{}) (string, error)
}
type NodeMeta struct {
	StructName string
	Prefix     string
	Middleware []string

	MethodName string
	HttpMethod string
	Path       string
}
