package router

////go:build ignore
//
//package main
//
//import (
//	"fmt"
//	"go/token"
//	"gtpkg"
//	"os"
//	"path/filepath"
//
//	"encoding/json"
//	"go/ast"
//	"go/parser"
//	"regexp"
//	"strings"
//
//	"github.com/bmatcuk/doublestar/v4"
//)
//
//func main() {
//	root := pkg.AppPath
//	mainPath := pkg.MainPath
//	outFilename := filepath.Join("./", root, "route_gen.go")
//	out, err := os.Create(outFilename)
//	if err != nil {
//		panic(err)
//	}
//	defer out.Close()
//
//	fs := os.DirFS(filepath.Join("./", root))
//
//	filenames, err := doublestar.Glob(fs, "*/controller/**.go")
//	fmt.Println(mainPath)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, filename := range filenames {
//		//file, err := fs.Open(filename)
//		//if err != nil {
//		//	fmt.Printf("打开文件失败 %s: %v\n", filename, err)
//		//	continue
//		//}
//		//defer file.Close()
//		//b, err := io.ReadAll(file)
//		//if err != nil {
//		//	fmt.Printf("读取文件失败 %s: %v\n", filename, err)
//		//	continue
//		//}
//		//
//		//fmt.Println(b)
//
//		p := NewParser()
//		p.ParseFile(filename)
//
//		//fset := token.NewFileSet()
//		//astFile, err := parser.ParseFile(fset, filename, b, parser.ParseComments)
//		//if err != nil {
//		//	fmt.Printf("解析 AST 失败 %s: %v\n", filename, err)
//		//	continue
//		//}
//	}
//	out.WriteString(`package app`)
//}
//
//var (
//	// @Router /path [method]
//	routerPattern = regexp.MustCompile(`@Router\s+(\S+)\s*\[(\w+)]`)
//
//	// @Security name
//	securityPattern = regexp.MustCompile(`@Security\s+(\w+)`)
//
//	// @Success code {type} schema "description"
//	successPattern = regexp.MustCompile(`@Success\s+(\d+)\s+{(\S+)}\s+(\S+)\s*"?([^"]*)"?`)
//
//	// @Param name paramType dataType required "description"
//	paramPattern = regexp.MustCompile(`@Param\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+"?([^"]*)"?`)
//
//	// @Summary text
//	summaryPattern = regexp.MustCompile(`@Summary\s+(.+)`)
//
//	// @Description text
//	descriptionPattern = regexp.MustCompile(`@Description\s+(.+)`)
//
//	// @Tags tag1,tag2
//	tagsPattern = regexp.MustCompile(`@Tags\s+(.+)`)
//
//	// @Accept mime
//	acceptPattern = regexp.MustCompile(`@Accept\s+(.+)`)
//
//	// @Produce mime
//	producePattern = regexp.MustCompile(`@Produce\s+(.+)`)
//
//	// @Title text
//	titlePattern = regexp.MustCompile(`@Title\s+(.+)`)
//
//	// @Version text
//	versionPattern = regexp.MustCompile(`@Version\s+(.+)`)
//
//	// @Host text
//	hostPattern = regexp.MustCompile(`@Host\s+(.+)`)
//
//	// @BasePath text
//	basePathPattern = regexp.MustCompile(`@BasePath\s+(.+)`)
//
//	// @SecurityDefinitions.apikey name
//	securityDefPattern = regexp.MustCompile(`@SecurityDefinitions\.(\w+)\s+(\w+)`)
//
//	// @in location
//	securityInPattern = regexp.MustCompile(`@in\s+(\w+)`)
//
//	// @name field
//	securityNamePattern = regexp.MustCompile(`@name\s+(\S+)`)
//)
//
//// Operation 表示一个 API 操作
//type Operation struct {
//	HTTPMethod  string                `json:"-"`
//	Path        string                `json:"-"`
//	Summary     string                `json:"summary,omitempty"`
//	Description string                `json:"description,omitempty"`
//	Tags        []string              `json:"tags,omitempty"`
//	Consumes    []string              `json:"consumes,omitempty"`
//	Produces    []string              `json:"produces,omitempty"`
//	Parameters  []*Parameter          `json:"parameters,omitempty"`
//	Responses   map[string]*Response  `json:"responses"`
//	Security    []map[string][]string `json:"security,omitempty"`
//}
//
//// Parameter 表示请求参数
//type Parameter struct {
//	Name        string  `json:"name"`
//	In          string  `json:"in"`
//	Description string  `json:"description,omitempty"`
//	Required    bool    `json:"required,omitempty"`
//	Type        string  `json:"type,omitempty"`
//	Schema      *Schema `json:"schema,omitempty"`
//}
//
//// Response 表示响应
//type Response struct {
//	Description string  `json:"description"`
//	Schema      *Schema `json:"schema,omitempty"`
//}
//
//// Schema 表示数据结构
//type Schema struct {
//	Type                 string             `json:"type,omitempty"`
//	Format               string             `json:"format,omitempty"`
//	Ref                  string             `json:"$ref,omitempty"`
//	Items                *Schema            `json:"items,omitempty"`
//	Properties           map[string]*Schema `json:"properties,omitempty"`
//	Required             []string           `json:"required,omitempty"`
//	Example              interface{}        `json:"example,omitempty"`
//	AdditionalProperties *Schema            `json:"additionalProperties,omitempty"`
//}
//
//// Info 表示 API 基本信息
//type Info struct {
//	Title          string `json:"title"`
//	Description    string `json:"description,omitempty"`
//	Version        string `json:"version"`
//	TermsOfService string `json:"termsOfService,omitempty"`
//}
//
//// SecurityScheme 表示安全认证方案
//type SecurityScheme struct {
//	Type        string `json:"type"`
//	Name        string `json:"name,omitempty"`
//	In          string `json:"in,omitempty"`
//	Description string `json:"description,omitempty"`
//	TokenURL    string `json:"tokenUrl,omitempty"`
//}
//
//// Specification 表示完整的 Swagger 规范
//type Specification struct {
//	Swagger             string                           `json:"swagger"`
//	Info                Info                             `json:"info"`
//	Host                string                           `json:"host,omitempty"`
//	BasePath            string                           `json:"basePath,omitempty"`
//	Paths               map[string]map[string]*Operation `json:"paths"`
//	Definitions         map[string]*Schema               `json:"definitions,omitempty"`
//	SecurityDefinitions map[string]*SecurityScheme       `json:"securityDefinitions,omitempty"`
//}
//
//// Parser Swagger 解析器
//type Parser struct {
//	spec          *Specification
//	currentOp     *Operation
//	definitions   map[string]*Schema
//	typeOverrides map[string]string
//}
//
//// NewParser 创建新的解析器
//func NewParser() *Parser {
//	return &Parser{
//		spec: &Specification{
//			Swagger: "2.0",
//			Info: Info{
//				Title:   "API Documentation",
//				Version: "1.0",
//			},
//			Paths:               make(map[string]map[string]*Operation),
//			Definitions:         make(map[string]*Schema),
//			SecurityDefinitions: make(map[string]*SecurityScheme),
//		},
//		definitions:   make(map[string]*Schema),
//		typeOverrides: make(map[string]string),
//	}
//}
//
//// ParseFile 解析单个 Go 文件
//func (p *Parser) ParseFile(filename string) error {
//	fset := token.NewFileSet()
//	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
//	if err != nil {
//		return fmt.Errorf("解析文件 %s 失败: %w", filename, err)
//	}
//
//	fmt.Println(file)
//
//	// 遍历 AST 节点
//	ast.Inspect(file, func(n ast.Node) bool {
//		switch x := n.(type) {
//		case *ast.Ast:
//			// 解析文件级别的注释（全局配置）
//			if x.Doc != nil {
//				p.parseGlobalAnnotations(x.Doc.Text())
//			}
//
//		case *ast.FuncDecl:
//			// 解析函数注释（API 操作）
//			if x.Doc != nil && x.Doc.List != nil {
//				p.parseFuncAnnotations(x)
//			}
//
//		case *ast.TypeSpec:
//			// 解析类型定义（结构体）
//			if x.Doc != nil {
//				p.parseTypeDefinition(x)
//			}
//		}
//		return true
//	})
//
//	return nil
//}
//
//// ParseDirectory 解析目录下的所有 Go 文件
//func (p *Parser) ParseDirectory(dir string) error {
//	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//
//		// 跳过隐藏目录和 vendor 目录
//		if info.IsDir() {
//			name := info.Name()
//			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
//				return filepath.SkipDir
//			}
//			return nil
//		}
//
//		// 只处理 .go 文件
//		if !strings.HasSuffix(path, ".go") {
//			return nil
//		}
//
//		// 跳过测试文件
//		if strings.HasSuffix(path, "_test.go") {
//			return nil
//		}
//
//		return p.ParseFile(path)
//	})
//}
//
//// parseGlobalAnnotations 解析全局注解
//func (p *Parser) parseGlobalAnnotations(comment string) {
//	lines := strings.Split(comment, "\n")
//	for _, line := range lines {
//		line = strings.TrimSpace(line)
//
//		if matches := titlePattern.FindStringSubmatch(line); matches != nil {
//			p.spec.Info.Title = matches[1]
//		} else if matches := versionPattern.FindStringSubmatch(line); matches != nil {
//			p.spec.Info.Version = matches[1]
//		} else if matches := descriptionPattern.FindStringSubmatch(line); matches != nil {
//			p.spec.Info.Description = matches[1]
//		} else if matches := hostPattern.FindStringSubmatch(line); matches != nil {
//			p.spec.Host = matches[1]
//		} else if matches := basePathPattern.FindStringSubmatch(line); matches != nil {
//			p.spec.BasePath = matches[1]
//		}
//	}
//}
//
//// parseFuncAnnotations 解析函数注解
//func (p *Parser) parseFuncAnnotations(funcDecl *ast.FuncDecl) {
//	comment := funcDecl.Doc.Text()
//
//	// 检查是否包含 @Router 注解
//	if !routerPattern.MatchString(comment) {
//		return
//	}
//
//	// 创建新的操作对象
//	op := &Operation{
//		Responses: make(map[string]*Response),
//	}
//
//	// 解析所有注解
//	p.parseOperationAnnotations(comment, op)
//
//	// 将操作添加到路径映射中
//	if op.Path != "" && op.HTTPMethod != "" {
//		pathKey := op.Path
//		if _, exists := p.spec.Paths[pathKey]; !exists {
//			p.spec.Paths[pathKey] = make(map[string]*Operation)
//		}
//		p.spec.Paths[pathKey][strings.ToLower(op.HTTPMethod)] = op
//	}
//}
//
//// parseOperationAnnotations 解析操作注解
//func (p *Parser) parseOperationAnnotations(comment string, op *Operation) {
//	lines := strings.Split(comment, "\n")
//	for _, line := range lines {
//		line = strings.TrimSpace(line)
//
//		// 匹配 @Router
//		if matches := routerPattern.FindStringSubmatch(line); matches != nil {
//			op.Path = matches[1]
//			op.HTTPMethod = strings.ToUpper(matches[2])
//		}
//
//		// 匹配 @Security
//		if matches := securityPattern.FindStringSubmatch(line); matches != nil {
//			op.Security = append(op.Security, map[string][]string{matches[1]: {}})
//		}
//
//		// 匹配 @Success
//		if matches := successPattern.FindStringSubmatch(line); matches != nil {
//			code := matches[1]
//			schemaType := matches[2]
//			schemaName := matches[3]
//			description := matches[4]
//
//			op.Responses[code] = &Response{
//				Description: description,
//				Schema:      p.parseSchema(schemaType, schemaName),
//			}
//		}
//
//		// 匹配 @Param
//		if matches := paramPattern.FindStringSubmatch(line); matches != nil {
//			param := &Parameter{
//				Name:        matches[1],
//				In:          matches[2],
//				Type:        matches[3],
//				Description: matches[5],
//			}
//
//			// 判断是否必需
//			required := strings.ToLower(matches[4])
//			param.Required = required == "true" || required == "required"
//
//			op.Parameters = append(op.Parameters, param)
//		}
//
//		// 匹配 @Summary
//		if matches := summaryPattern.FindStringSubmatch(line); matches != nil {
//			op.Summary = matches[1]
//		}
//
//		// 匹配 @Description
//		if matches := descriptionPattern.FindStringSubmatch(line); matches != nil {
//			op.Description = matches[1]
//		}
//
//		// 匹配 @Tags
//		if matches := tagsPattern.FindStringSubmatch(line); matches != nil {
//			tags := strings.Split(matches[1], ",")
//			for _, tag := range tags {
//				tag = strings.TrimSpace(tag)
//				if tag != "" {
//					op.Tags = append(op.Tags, tag)
//				}
//			}
//		}
//
//		// 匹配 @Accept
//		if matches := acceptPattern.FindStringSubmatch(line); matches != nil {
//			op.Consumes = append(op.Consumes, matches[1])
//		}
//
//		// 匹配 @Produce
//		if matches := producePattern.FindStringSubmatch(line); matches != nil {
//			op.Produces = append(op.Produces, matches[1])
//		}
//	}
//}
//
//// parseSchema 解析 Schema 类型
//func (p *Parser) parseSchema(schemaType string, schemaName string) *Schema {
//	schema := &Schema{}
//
//	switch schemaType {
//	case "object":
//		// 引用类型
//		schema.Ref = "#/definitions/" + schemaName
//	case "array":
//		schema.Type = "array"
//		schema.Items = &Schema{
//			Ref: "#/definitions/" + schemaName,
//		}
//	case "string":
//		schema.Type = "string"
//	case "integer":
//		schema.Type = "integer"
//	case "number":
//		schema.Type = "number"
//	case "boolean":
//		schema.Type = "boolean"
//	default:
//		schema.Type = schemaType
//	}
//
//	return schema
//}
//
//// parseTypeDefinition 解析类型定义
//func (p *Parser) parseTypeDefinition(typeSpec *ast.TypeSpec) {
//	// 只处理结构体类型
//	structType, ok := typeSpec.Type.(*ast.StructType)
//	if !ok {
//		return
//	}
//
//	typeName := typeSpec.Name.Name
//	schema := &Schema{
//		Type:       "object",
//		Properties: make(map[string]*Schema),
//	}
//
//	// 解析结构体字段
//	for _, field := range structType.Fields.List {
//		if len(field.Names) == 0 {
//			continue
//		}
//
//		jsonTag := ""
//
//		// 提取 json tag
//		if field.Tag != nil {
//			tagValue := field.Tag.Value
//			if strings.Contains(tagValue, "json:") {
//				tagParts := strings.Split(tagValue, `"`)
//				for i, part := range tagParts {
//					if part == "json:" && i+1 < len(tagParts) {
//						jsonTag = tagParts[i+1]
//						break
//					}
//				}
//			}
//		}
//
//		// 如果没有 json tag，使用字段名
//		if jsonTag == "" || jsonTag == "-" {
//			continue
//		}
//
//		// 解析字段类型
//		fieldSchema := p.parseFieldType(field.Type)
//
//		// 提取 example
//		if field.Tag != nil {
//			tagValue := field.Tag.Value
//			if strings.Contains(tagValue, "example:") {
//				example := p.extractTagValue(tagValue, "example")
//				if example != "" {
//					fieldSchema.Example = example
//				}
//			}
//		}
//
//		schema.Properties[jsonTag] = fieldSchema
//	}
//
//	// 存储定义
//	p.definitions[typeName] = schema
//	p.spec.Definitions[typeName] = schema
//}
//
//// parseFieldType 解析字段类型
//func (p *Parser) parseFieldType(expr ast.Expr) *Schema {
//	switch t := expr.(type) {
//	case *ast.Ident:
//		// 基本类型
//		return p.parseBasicType(t.Name)
//
//	case *ast.StarExpr:
//		// 指针类型，递归解析
//		return p.parseFieldType(t.X)
//
//	case *ast.ArrayType:
//		// 数组类型
//		return &Schema{
//			Type:  "array",
//			Items: p.parseFieldType(t.Elt),
//		}
//
//	case *ast.MapType:
//		// Map 类型
//		return &Schema{
//			Type:                 "object",
//			AdditionalProperties: &Schema{},
//		}
//
//	case *ast.SelectorExpr:
//		// 包限定类型
//		if ident, ok := t.X.(*ast.Ident); ok {
//			typeName := ident.Name + "." + t.Sel.Name
//			return &Schema{
//				Ref: "#/definitions/" + typeName,
//			}
//		}
//	}
//
//	return &Schema{Type: "object"}
//}
//
//// parseBasicType 解析基本类型
//func (p *Parser) parseBasicType(typeName string) *Schema {
//	switch typeName {
//	case "string":
//		return &Schema{Type: "string"}
//	case "int", "int8", "int16", "int32", "int64",
//		"uint", "uint8", "uint16", "uint32", "uint64":
//		return &Schema{Type: "integer"}
//	case "float32", "float64":
//		return &Schema{Type: "number", Format: "double"}
//	case "bool":
//		return &Schema{Type: "boolean"}
//	case "time.Time":
//		return &Schema{Type: "string", Format: "date-time"}
//	default:
//		// 自定义类型，尝试引用
//		return &Schema{
//			Ref: "#/definitions/" + typeName,
//		}
//	}
//}
//
//// extractTagValue 从标签字符串中提取值
//func (p *Parser) extractTagValue(tag, key string) string {
//	parts := strings.Split(tag, `"`)
//	for i, part := range parts {
//		if strings.HasPrefix(part, key+":") && i+1 < len(parts) {
//			return parts[i+1]
//		}
//	}
//	return ""
//}
//
//// Generate 生成 Swagger 文档
//func (p *Parser) Generate(outputDir string) error {
//	// 确保输出目录存在
//	if err := os.MkdirAll(outputDir, 0755); err != nil {
//		return fmt.Errorf("创建输出目录失败: %w", err)
//	}
//
//	// 生成 swagger.json
//	jsonPath := filepath.Join(outputDir, "swagger.json")
//	jsonData, err := json.MarshalIndent(p.spec, "", "  ")
//	if err != nil {
//		return fmt.Errorf("序列化 JSON 失败: %w", err)
//	}
//
//	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
//		return fmt.Errorf("写入 JSON 文件失败: %w", err)
//	}
//
//	// 生成 swagger.yaml
//	yamlPath := filepath.Join(outputDir, "swagger.yaml")
//	yamlData, err := p.generateYAML()
//	if err != nil {
//		return fmt.Errorf("生成 YAML 失败: %w", err)
//	}
//
//	if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
//		return fmt.Errorf("写入 YAML 文件失败: %w", err)
//	}
//
//	fmt.Printf("Swagger 文档已生成到 %s\n", outputDir)
//	fmt.Printf("  - %s\n", jsonPath)
//	fmt.Printf("  - %s\n", yamlPath)
//
//	return nil
//}
//
//// generateYAML 生成 YAML 格式（简化版）
//func (p *Parser) generateYAML() ([]byte, error) {
//	// 这里可以使用 yaml.Marshal，为了简化先返回 JSON
//	// 实际项目中建议引入 gopkg.in/yaml.v3
//	return json.MarshalIndent(p.spec, "", "  ")
//}
//
//// GetSpecification 获取规范对象
//func (p *Parser) GetSpecification() *Specification {
//	return p.spec
//}
//
//// AddSecurityDefinition 添加安全认证定义
//func (p *Parser) AddSecurityDefinition(name string, scheme *SecurityScheme) {
//	p.spec.SecurityDefinitions[name] = scheme
//}
//
//// ParseSecurityDefinitions 解析安全认证定义注解
//func (p *Parser) ParseSecurityDefinitions(comment string) {
//	lines := strings.Split(comment, "\n")
//
//	var currentScheme *SecurityScheme
//	var currentName string
//
//	for _, line := range lines {
//		line = strings.TrimSpace(line)
//
//		// 匹配 @SecurityDefinitions.apikey Name
//		if matches := securityDefPattern.FindStringSubmatch(line); matches != nil {
//			if currentName != "" && currentScheme != nil {
//				p.AddSecurityDefinition(currentName, currentScheme)
//			}
//
//			currentName = matches[2]
//			currentScheme = &SecurityScheme{
//				Type: matches[1],
//			}
//		}
//
//		// 匹配 @in header/query/cookie
//		if matches := securityInPattern.FindStringSubmatch(line); matches != nil {
//			if currentScheme != nil {
//				currentScheme.In = matches[1]
//			}
//		}
//
//		// 匹配 @name field
//		if matches := securityNamePattern.FindStringSubmatch(line); matches != nil {
//			if currentScheme != nil {
//				currentScheme.Name = matches[1]
//			}
//		}
//
//		// 匹配 @description
//		if strings.HasPrefix(line, "@description") || strings.HasPrefix(line, "// @description") {
//			if currentScheme != nil {
//				desc := strings.TrimPrefix(line, "@description")
//				desc = strings.TrimPrefix(desc, "// @description")
//				currentScheme.Description = strings.TrimSpace(desc)
//			}
//		}
//	}
//
//	// 保存最后一个方案
//	if currentName != "" && currentScheme != nil {
//		p.AddSecurityDefinition(currentName, currentScheme)
//	}
//}
