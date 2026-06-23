package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gta-spec/utils/slog"
)

var Logger *slog.Logger
var rollerWriter *_slog.Roller

func init() {
	rollerWriter, err := _slog.NewRoller(
		"runtime/log/%Y-%m-%d/%Y-%m-%d.log",
		200*1024*1024, // 最大两百兆
		&_slog.Options{
			MaxBackups: 10,
			MaxAge:     30 * 24 * time.Hour, // 30 days
			LocalTime:  true,
			Compress:   false,
		},
	)
	if err != nil {
		panic(err)
	}
	// 日志双写
	writers := []io.Writer{rollerWriter, os.Stdout}

	Logger = New(
		WithLevel(slog.LevelDebug),
		WithCallerSkip(1),
		WithWriter(writers...),
		WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
			//if a.Key == slog.LevelKey {
			//	if a.Value.Any().(slog.Level) == slog.LevelError {
			//	}
			//}
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.DateTime + ".000"))
			}
			return a
		}),
	)
	slog.SetDefault(New(WithLevel(slog.LevelDebug), WithWriter(writers...)))
}

func Close() {
	if rollerWriter != nil {
		rollerWriter.Close()
	}
}

// Record 记录调试信息
// 参数:
//
//	msg: 调试信息
//	type: 信息类型 info
func Record(ctx context.Context, level slog.Level, msg string, args ...any) {
	//slog.Log(ctx, level, msg, args...)
	Write(ctx, level, msg, args...)
}
func RecordAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	//slog.LogAttrs(ctx, level, msg, attrs...)
	WriteAttrs(ctx, level, msg, attrs...)
}

// Save 把保存在内存中的日志信息写入
func Save() {

}

// Write 把保存在内存中的日志信息写入
func Write(ctx context.Context, level slog.Level, msg string, args ...any) {
	Logger.Log(ctx, level, msg, args...)
}

// WriteAttrs 把保存在内存中的日志信息写入
func WriteAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	Logger.LogAttrs(ctx, level, msg, attrs...)
}
