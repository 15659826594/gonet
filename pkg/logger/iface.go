package logger

import (
	"context"
)

type FieldType uint8

// Field 日志字段（键值对）
type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

// Logger 通用日志接口（PSR-3 兼容）
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Notice(msg string, fields ...Field)
	Warning(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Critical(msg string, fields ...Field)
	Alert(msg string, fields ...Field)
	Emergency(msg string, fields ...Field)

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Noticef(template string, args ...interface{})
	Warningf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Criticalf(template string, args ...interface{})
	Alertf(template string, args ...interface{})
	Emergencyf(template string, args ...interface{})

	// DebugCtx 带上下文的日志（可选扩展）
	DebugCtx(ctx context.Context, msg string, fields ...Field)
	InfoCtx(ctx context.Context, msg string, fields ...Field)
	NoticeCtx(ctx context.Context, msg string, fields ...Field)
	WarningCtx(ctx context.Context, msg string, fields ...Field)
	ErrorCtx(ctx context.Context, msg string, fields ...Field)
	CriticalCtx(ctx context.Context, msg string, fields ...Field)
	AlertCtx(ctx context.Context, msg string, fields ...Field)
	EmergencyCtx(ctx context.Context, msg string, fields ...Field)

	Sync() error
}
