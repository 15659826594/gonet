package utils

import (
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/language"
)

// TrimExt 移除文件后缀名
func TrimExt(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// StrShuffle str_shuffle() 随机地打乱字符串中的所有字符。
func StrShuffle(str string) string {
	runes := []rune(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]rune, len(runes))
	for i, v := range r.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// Substr substr() 返回字符串的子串
func Substr(str string, start uint, length int) string {
	if length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

// Langset 获取当前的语言
func Langset(header string) language.Tag {
	tags, _, _ := language.ParseAcceptLanguage(header)
	matcher := language.NewMatcher(tags)
	matched, _, _ := matcher.Match()
	return matched
}
