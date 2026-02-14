package utils

import (
	"net/url"
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
