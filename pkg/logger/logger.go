package logger

import (
	"context"
	"log/slog"
)

var log *slog.Logger

func init() {
	def := NewLogger(&Option{
		callerSkip:  0,
		writer:      nil,
		AddSource:   false,
		Level:       nil,
		ReplaceAttr: nil,
	})
	slog.SetDefault(def)
}

// Record 记录调试信息
// 参数:
//
//	msg: 调试信息
//	type: 信息类型 info
func Record(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}
func RecordAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, level, msg, attrs...)
}

// Save 把保存在内存中的日志信息写入
func Save() {

}

// Write 把保存在内存中的日志信息写入
func Write(ctx context.Context, level slog.Level, msg string, args ...any) {
	log.Log(ctx, level, msg, args...)
}
