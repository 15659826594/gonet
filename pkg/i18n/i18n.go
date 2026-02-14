package i18n

import (
	"encoding/json"
	"fmt"
	"gonet/pkg"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//fmt.Println(i18n.T("admin/user/rule", "Menu tips"))
//fmt.Println(i18n.T("admin/dashboard", "Total user"))
//fmt.Println(i18n.T("index/user", "Maxsuccessions"))
//fmt.Println(i18n.T("index", "Forgot password"))
//fmt.Println(i18n.T("", "PSR-4 error"))

func init() {
	_, filename, _, _ := runtime.Caller(0)
	root = &Translate{
		Path:   "",
		Bundle: i18n.NewBundle(defaultLang),
	}

	root.LoadDir(filepath.Join(filepath.Dir(filename), "lang"))

	matches, _ := filepath.Glob(pkg.APP_PATH + fmt.Sprintf("*%slang", pkg.DS))
	for _, match := range matches {
		moduleName := strings.Split(match, pkg.DS)[1]
		module := &Translate{
			Path:   moduleName,
			Bundle: i18n.NewBundle(defaultLang),
		}
		// 然后遍历目录
		files, _ := os.ReadDir(match)
		for _, file := range files {
			if file.IsDir() { // 判断是否为文件夹
				if tag, err := language.Parse(file.Name()); err == nil {
					module.LoadDirWithTag(filepath.Join(match, file.Name()), tag)
				}
			}
		}

		module.LoadDir(match)
		root.Children = append(root.Children, module)
	}
	root.InitLocalizer()
}

var (
	defaultLang = language.English
	root        *Translate
)

var Weights = map[string]int{
	"zh-CN": 10, // 更高权重
	"en":    5,  // 较低权重
}

type Translate struct {
	Path      string
	Bundle    *i18n.Bundle
	Langs     []string
	Localizer *i18n.Localizer
	Children  []*Translate
}

func (t *Translate) LoadDir(dir string) {
	files, _ := os.ReadDir(dir)
	// 遍历文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(dir, file.Name())
		ext := filepath.Ext(path)
		switch ext {
		case ".json":
			t.Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
			if tag, err := language.Parse(strings.TrimSuffix(file.Name(), ext)); err == nil {
				t.Bundle.LoadMessageFile(path)
				t.Langs = append(t.Langs, tag.String())
			}
		}
	}
}

func (t *Translate) LoadDirWithTag(dir string, tag language.Tag) {
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		filename := strings.TrimSuffix(file.Name(), ext)
		child := t.GetNode(filename)
		if child == nil {
			child = &Translate{
				Path:   filename,
				Bundle: i18n.NewBundle(defaultLang),
			}
			child.Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
			t.Children = append(t.Children, child)
		}
		switch ext {
		case ".json":
			if jsonData, err := os.ReadFile(filepath.Join(dir, file.Name())); err == nil {
				if _, err := child.Bundle.ParseMessageFileBytes(jsonData, filepath.Join(dir, FormatFilename(file.Name(), tag))); err == nil {
					child.Langs = append(child.Langs, tag.String())
				}
			}
		case "":
			child.LoadDirWithTag(filepath.Join(dir, filename), tag)
		}
	}
}

// GetNode 根据路径 获取子节点
func (t *Translate) GetNode(path string) *Translate {
	for _, child := range t.Children {
		if child.Path == path {
			return child
		}
	}
	return nil
}

// GetTransLink 获取翻译链路
func (t *Translate) GetTransLink(paths []string) []*Translate {
	trans := []*Translate{t}
	tmp := t
	for _, path := range paths {
		result := tmp.GetNode(path)
		if result == nil {
			break
		}
		trans = append(trans, result)
		tmp = result
	}
	return trans
}

// SortLangs 语言权种排序
func (t *Translate) SortLangs() []string {
	langs := make([]string, len(t.Langs))
	copy(langs, t.Langs)

	for i := 0; i < len(langs)-1; i++ {
		for j := i + 1; j < len(langs); j++ {
			if Weights[langs[i]] < Weights[langs[j]] {
				langs[i], langs[j] = langs[j], langs[i]
			}
		}
	}
	return langs
}

func (t *Translate) InitLocalizer() {
	// 为当前节点创建本地化器，使用排序后的语言列表
	t.Localizer = i18n.NewLocalizer(t.Bundle, t.SortLangs()...)

	// 递归为所有子节点初始化本地化器
	for _, child := range t.Children {
		child.InitLocalizer()
	}
}

// FormatFilename 格式化文件名，将原始文件名转换为包含语言标签的格式
// 例如: index.json + zh-CN -> index.zh-cn.json
// 参数:
//   - filename: 原始文件名
//   - tag: 语言标签
//
// 返回值: 添加了语言标签的文件名
func FormatFilename(filename string, tag language.Tag) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	langTag := strings.ToLower(tag.String())
	return fmt.Sprintf("%s.%s%s", nameWithoutExt, langTag, ext)
}

// T 根据URL路径、消息ID和模板参数进行国际化翻译
// url: 模块/控制器路径，用于定位翻译文件位置
// messageID: 要翻译的字符串标识符
// templates: 可选的模板数据，用于替换翻译文本中的占位符
// 返回值: 翻译后的字符串，如果未找到对应翻译则返回messageID
func T(url string, messageID string, templates ...map[string]any) string {
	trans := root.GetTransLink(strings.Split(url, "/"))

	// 从链路末尾开始向前查找翻译
	for i := len(trans) - 1; i >= 0; i-- {
		node := trans[i]
		if node != nil && node.Localizer != nil {
			// 尝试翻译
			var config *i18n.LocalizeConfig
			if len(templates) > 0 {
				config = &i18n.LocalizeConfig{
					MessageID:    messageID,
					TemplateData: templates[0],
				}
			} else {
				config = &i18n.LocalizeConfig{
					MessageID: messageID,
				}
			}

			if result, err := node.Localizer.Localize(config); err == nil {
				return result
			}
		}
	}
	return messageID
}
