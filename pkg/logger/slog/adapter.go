package slog

import (
	"context"
	"fmt"
	"gonet/pkg/logger"
	"gonet/pkg/logger/level"
	"log/slog"
	"time"
)

// Logger slog 适配器
type Logger struct {
	logger *slog.Logger
}

// Adapter 创建 slog 日志适配器
func Adapter(handler slog.Handler) logger.Logger {
	return &Logger{logger: slog.New(handler)}
}

// 将 Field 转换为 slog.Attr
func toFields(fields []logger.Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for _, f := range fields {
		attrs = append(attrs, slog.Any(f.Key, f.Interface))
	}
	return attrs
}

func toLevel(level level.Level) slog.Level {
	return slog.Level(level)
}

// 实现 PSR-3 方法

func (s *Logger) Debug(msg string, fields ...logger.Field) {
	s.log(level.Debug, msg, fields...)
}

func (s *Logger) Notice(msg string, fields ...logger.Field) {
	s.log(level.Notice, msg, fields...)
}

func (s *Logger) Info(msg string, fields ...logger.Field) {
	s.log(level.Info, msg, fields...)
}

func (s *Logger) Warning(msg string, fields ...logger.Field) {
	s.log(level.Warning, msg, fields...)
}

func (s *Logger) Error(msg string, fields ...logger.Field) {
	s.log(level.Error, msg, fields...)
}

func (s *Logger) Critical(msg string, fields ...logger.Field) {
	s.log(level.Critical, msg, fields...)
}

func (s *Logger) Alert(msg string, fields ...logger.Field) {
	s.log(level.Alert, msg, fields...)
}

func (s *Logger) Emergency(msg string, fields ...logger.Field) {
	s.log(level.Emergency, msg, fields...)
}

// 传统格式

func (s *Logger) Debugf(template string, args ...interface{}) {
	s.logf(level.Debug, template, args...)
}

func (s *Logger) Infof(template string, args ...interface{}) {
	s.logf(level.Info, template, args...)
}

func (s *Logger) Noticef(template string, args ...interface{}) {
	s.logf(level.Notice, template, args...)
}

func (s *Logger) Warningf(template string, args ...interface{}) {
	s.logf(level.Warning, template, args...)
}

func (s *Logger) Errorf(template string, args ...interface{}) {
	s.logf(level.Error, template, args...)
}

func (s *Logger) Criticalf(template string, args ...interface{}) {
	s.logf(level.Critical, template, args...)
}

func (s *Logger) Alertf(template string, args ...interface{}) {
	s.logf(level.Alert, template, args...)
}

func (s *Logger) Emergencyf(template string, args ...interface{}) {
	s.logf(level.Emergency, template, args...)
}

// 带上下文的日志方法

func (s *Logger) DebugCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Debug, msg, fields...)
}

func (s *Logger) InfoCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Info, msg, fields...)
}

func (s *Logger) NoticeCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Notice, msg, fields...)
}

func (s *Logger) WarningCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Warning, msg, fields...)
}

func (s *Logger) ErrorCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Error, msg, fields...)
}

func (s *Logger) CriticalCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Critical, msg, fields...)
}

func (s *Logger) AlertCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Alert, msg, fields...)
}

func (s *Logger) EmergencyCtx(ctx context.Context, msg string, fields ...logger.Field) {
	s.logCtx(ctx, level.Emergency, msg, fields...)
}

// log 统一日志处理
func (s *Logger) log(level level.Level, msg string, fields ...logger.Field) {
	// 构造 slog.Record
	record := slog.NewRecord(time.Now(), toLevel(level), msg, 0)
	record.AddAttrs(toFields(fields)...)
	_ = s.logger.Handler().Handle(context.Background(), record)
}

func (s *Logger) logf(level level.Level, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	// 构造 slog.Record
	record := slog.NewRecord(time.Now(), toLevel(level), formattedMsg, 0)
	_ = s.logger.Handler().Handle(context.Background(), record)
}

// logCtx 统一带上下文的日志处理
func (s *Logger) logCtx(ctx context.Context, level level.Level, msg string, fields ...logger.Field) {
	// 构造 slog.Record
	record := slog.NewRecord(time.Now(), toLevel(level), msg, 0)
	record.AddAttrs(toFields(fields)...)
	_ = s.logger.Handler().Handle(ctx, record)
}

func (s *Logger) Sync() error {
	return nil
}
