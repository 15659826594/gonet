package utils

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"
)

// Url 结构体用于处理URL相关功能
type Url struct {
	root string
}

// Build 构建URL地址
func (u *Url) Build(baseUrl, target string, vars string) string {
	base, _ := url.Parse(baseUrl)
	relative, _ := url.Parse(target)
	newUrl := base.ResolveReference(relative)

	params, _ := url.ParseQuery(base.RawQuery)
	parsed, err := url.ParseQuery(vars)
	if err == nil {
		for key, val := range parsed {
			if len(val) > 0 {
				params[key] = val
			}
		}
	}
	newUrl.RawQuery = params.Encode()

	return newUrl.String()
}
func Stack(skip int) string {
	var buf strings.Builder

	// +1 是为了跳过 Stack 函数本身
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 获取函数名
		fn := runtime.FuncForPC(pc)
		funcName := "???"
		if fn != nil {
			funcName = fn.Name()
			// 简化函数名，去掉包路径
			if idx := strings.LastIndex(funcName, "/"); idx != -1 {
				funcName = funcName[idx+1:]
			}
		}

		// 简化文件路径，只保留最后两部分
		shortFile := file
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			if idx2 := strings.LastIndex(file[:idx], "/"); idx2 != -1 {
				shortFile = file[idx2+1:]
			}
		}

		buf.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", funcName, shortFile, line))
	}

	return buf.String()
}
