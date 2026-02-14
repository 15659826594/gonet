package level

import (
	"bytes"
	"errors"
	"fmt"
)

// 日志级别说明
// Emergency  系统不可用        数据库崩溃
// Alert      必须立即处理      安全漏洞
// Critical   严重错误          应用组件失效
// Error      运行时错误        文件写入失败
// Warning    警告              弃用 API 调用
// Notice     正常但值得注意    用户登录
// Info       信息性消息        API 请求记录
// Debug      调试信息          开发阶段详细日志

type Level int

// PSR-3 标准日志级别
const (
	Debug Level = iota*4 - 4
	Info
	Warning
	Error
	Critical
	Alert
	Emergency
	Notice Level = 2 // 特殊定义，保持原有值
)

// ParseLevel 解析字符串表示的日志级别
func ParseLevel(text string) (Level, error) {
	var level Level
	err := level.UnmarshalText([]byte(text))
	return level, err
}

// String 返回小写ASCII格式的日志级别表示
func (l Level) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warning:
		return "warn"
	case Error:
		return "error"
	case Critical:
		return "critical"
	case Alert:
		return "alert"
	case Emergency:
		return "emergency"
	case Notice:
		return "notice"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// CapitalString 返回大写ASCII格式的日志级别表示
func (l Level) CapitalString() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARN"
	case Error:
		return "ERROR"
	case Critical:
		return "CRITICAL"
	case Alert:
		return "ALERT"
	case Emergency:
		return "EMERGENCY"
	case Notice:
		return "NOTICE"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// MarshalText 将 Level 编组为文本
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText 从文本解组 Level
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errors.New("can't unmarshal a nil *Level")
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = Debug
	case "info", "INFO", "":
		*l = Info
	case "warn", "WARN":
		*l = Warning
	case "error", "ERROR":
		*l = Error
	case "critical", "CRITICAL":
		*l = Critical
	case "alert", "ALERT":
		*l = Alert
	case "emergency", "EMERGENCY":
		*l = Emergency
	case "notice", "NOTICE":
		*l = Notice
	default:
		return false
	}
	return true
}

// Set 实现 flag.Value 接口
func (l *Level) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

// Get 实现 flag.Getter 接口
func (l *Level) Get() interface{} {
	return *l
}

// Enabled 判断给定级别是否启用
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// LevelEnabler 接口决定指定日志级别是否在记录消息时启用
type LevelEnabler interface {
	Enabled(Level) bool
}
