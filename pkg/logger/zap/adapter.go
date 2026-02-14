package zap

import (
	"context"
	"gonet/pkg/logger"
	"gonet/pkg/logger/level"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//const (
//	TraceLevel  zapcore.Level = -2 // 比Debug更详细的追踪级别
//	NoticeLevel zapcore.Level = 0  // 这里为了演示，若需避免与Info冲突，可调整默认级别映射，或选其他未占用数值（如6）
//)
//
//// 第二步：初始化自定义级别（注册级别名称，让日志输出显示可读名称而非数值）
//func init() {
//	// 注册自定义级别名称：第一个参数是级别数值，第二个参数是输出时显示的名称
//	zapcore.RegisterLevelName(int8(TraceLevel), "TRACE")
//	zapcore.RegisterLevelName(int8(NoticeLevel), "NOTICE")
//	zapcore.NewIncreaseLevelCore()
//}

// Logger zap 适配器
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// Adapter 创建 zap 日志适配器
func Adapter(logger *zap.Logger) logger.Logger {
	return &Logger{logger: logger, sugar: logger.Sugar()}
}

// 将 Field 转换为 zap.Field
func toFields(fields ...logger.Field) []zap.Field {
	fs := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		fs = append(fs, zap.Field{
			Key:       f.Key,
			Type:      zapcore.FieldType(f.Type),
			Integer:   f.Integer,
			String:    f.String,
			Interface: f.Interface,
		})
	}
	return fs
}

func toLevel(lev level.Level) zapcore.Level {
	switch lev {
	case level.Debug:
		return zapcore.DebugLevel
	case level.Info:
		return zapcore.InfoLevel
	case level.Notice:
		return zapcore.InfoLevel
	case level.Warning:
		return zapcore.WarnLevel
	case level.Error:
		return zapcore.ErrorLevel
	case level.Critical:
		return zapcore.DPanicLevel
	case level.Alert:
		return zapcore.PanicLevel
	case level.Emergency:
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

// 实现 PSR-3 方法

func (z *Logger) Debug(msg string, fields ...logger.Field) {
	z.log(level.Debug, msg, fields...)
}

func (z *Logger) Info(msg string, fields ...logger.Field) {
	z.log(level.Info, msg, fields...)
}

func (z *Logger) Notice(msg string, fields ...logger.Field) {
	z.log(level.Notice, msg, fields...)
}

func (z *Logger) Warning(msg string, fields ...logger.Field) {
	z.log(level.Warning, msg, fields...)
}

func (z *Logger) Error(msg string, fields ...logger.Field) {
	z.log(level.Error, msg, fields...)
}

func (z *Logger) Critical(msg string, fields ...logger.Field) {
	z.log(level.Critical, msg, fields...) // zap 的 Critical 对应 DPanicLevel
}

func (z *Logger) Alert(msg string, fields ...logger.Field) {
	z.log(level.Alert, msg, fields...)
}

func (z *Logger) Emergency(msg string, fields ...logger.Field) {
	z.log(level.Emergency, msg, fields...)
}

// 传统格式

func (z *Logger) Debugf(template string, args ...interface{}) {
	z.logf(level.Debug, template, args...)
}

func (z *Logger) Infof(template string, args ...interface{}) {
	z.logf(level.Info, template, args...)
}

func (z *Logger) Noticef(template string, args ...interface{}) {
	z.logf(level.Notice, template, args...)
}

func (z *Logger) Warningf(template string, args ...interface{}) {
	z.logf(level.Warning, template, args...)
}

func (z *Logger) Errorf(template string, args ...interface{}) {
	z.logf(level.Error, template, args...)
}

func (z *Logger) Criticalf(template string, args ...interface{}) {
	z.logf(level.Critical, template, args...)
}

func (z *Logger) Alertf(template string, args ...interface{}) {
	z.logf(level.Alert, template, args...)
}

func (z *Logger) Emergencyf(template string, args ...interface{}) {
	z.logf(level.Emergency, template, args...)
}

// 带上下文的日志方法

func (z *Logger) DebugCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Debug, msg, fields...)
}

func (z *Logger) InfoCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Info, msg, fields...)
}

func (z *Logger) NoticeCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Notice, msg, fields...)
}

func (z *Logger) WarningCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Warning, msg, fields...)
}

func (z *Logger) ErrorCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Error, msg, fields...)
}

func (z *Logger) CriticalCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Critical, msg, fields...)
}

func (z *Logger) AlertCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Alert, msg, fields...)
}

func (z *Logger) EmergencyCtx(ctx context.Context, msg string, fields ...logger.Field) {
	z.logCtx(ctx, level.Emergency, msg, fields...)
}

// log 统一日志处理
func (z *Logger) log(level level.Level, template string, fields ...logger.Field) {
	z.logger.WithOptions(zap.AddCallerSkip(1)).Check(toLevel(level), template).Write(toFields(fields...)...)
}

// logf 传统日志处理
func (z *Logger) logf(level level.Level, template string, args ...interface{}) {
	z.sugar.WithOptions(zap.AddCallerSkip(1)).Logf(toLevel(level), template, args...)
}

// logCtx 统一带上下文的日志处理
func (z *Logger) logCtx(ctx context.Context, level level.Level, template string, fields ...logger.Field) {
	zapFields := toFields(fields...)
	if ctx != nil {
		// 可从 ctx 中提取信息加入日志，如 traceID
		if traceId, ok := ctx.Value("traceId").(string); ok {
			zapFields = append(zapFields, zap.String("traceId", traceId))
		}
	}
	z.logger.WithOptions(zap.AddCallerSkip(1)).Check(toLevel(level), template).Write(zapFields...)
}

func (z *Logger) Sync() error {
	return z.logger.Sync()
}
