package utils

import (
	"reflect"
	"slices"
)

// IsNil 检查给定的值是否为nil
// 支持指针、切片、映射、通道、函数和接口类型
func IsNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

// Ternary 三目运算函数
// condition: 条件表达式
// trueValue: 条件为真时返回的值
// falseValue: 条件为假时返回的值
func Ternary[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// FirstNonNil (泛型)查找第一个非零值
func FirstNonNil[T comparable](value T, defaultValues ...T) T {
	var zero T
	if value != zero {
		return value
	}

	// 遍历默认值列表，返回第一个非零值
	for _, defaultValue := range defaultValues {
		if defaultValue != zero {
			return defaultValue
		}
	}

	// 如果没有找到非零默认值，则返回类型的零值
	return zero
}

// GetStructProperty 传入一个对象 通过反射获取属性
func GetStructProperty(object any, props ...string) map[string]any {
	t := reflect.TypeOf(object)
	v := reflect.ValueOf(object)

	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	ret := make(map[string]any)

	// 如果没有指定属性，则返回所有字段
	if len(props) == 0 {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			ret[field.Name] = value.Interface()
		}
		return ret
	}

	// 查找指定的字段
	for _, prop := range props {
		if field := v.FieldByName(prop); field.IsValid() {
			ret[prop] = field.Interface()
		}
	}

	return ret
}

// GetStructMethods 传入一个对象 通过反射获取属性, 如果没有指定属性名则返回所有方法
func GetStructMethods(object any, names ...string) map[string]any {
	v := reflect.ValueOf(object)

	ret := make(map[string]any)
	numMethods := v.NumMethod()

	for i := 0; i < numMethods; i++ {
		method := v.Type().Method(i)
		name := method.Name
		if len(names) == 0 || slices.Contains(names, name) {
			ret[name] = v.Method(i).Interface()
		}
	}

	return ret
}

// CaseSnake 驼峰转下划线命名
func CaseSnake(s string) string {
	if len(s) == 0 {
		return ""
	}

	buf := make([]byte, 0, len(s)*2)
	prevLower := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 && prevLower {
				buf = append(buf, '_')
			}
			buf = append(buf, c+32) // ASCII码转小写
			prevLower = false
		} else {
			if c == '_' && i > 0 && i < len(s)-1 {
				next := s[i+1]
				if next >= 'a' && next <= 'z' {
					next -= 32
				}
			}
			prevLower = (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
			buf = append(buf, c)
		}
	}
	return string(buf)
}

// SnakeCase 下划线转驼峰命名(upperFirst 大驼峰)
func SnakeCase(s string, args ...bool) string {
	upperFirst := false
	if len(args) > 0 {
		upperFirst = args[0]
	}
	buf := make([]byte, 0, len(s))
	nextUpper := upperFirst

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '_':
			nextUpper = true
		case nextUpper:
			if i == 0 && c >= 'A' && c <= 'Z' {
				c += 32
			} else if c >= 'a' && c <= 'z' {
				c -= 32 // ASCII码转大写
			}
			buf = append(buf, c)
			nextUpper = false
		default:
			buf = append(buf, c)
		}
	}
	return string(buf)
}
