package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// LogHandler 模拟 Logback 的 PatternLayout
type LogHandler struct {
	slog.Handler
	writer     io.Writer
	callerSkip int
	addSource  bool
}

type Option struct {
	callerSkip  int
	writer      io.Writer
	AddSource   bool
	Level       slog.Leveler
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func NewLogger(option *Option) *slog.Logger {
	return slog.New(&LogHandler{
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       option.Level,
			ReplaceAttr: option.ReplaceAttr,
		}),
		addSource:  option.AddSource,
		callerSkip: option.callerSkip,
	})
}

// Handle 核心格式化逻辑
func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	// 1. %d{yy-MM-dd HH:mm:ss.SSS}
	timeStr := r.Time.Format(time.DateTime + ".000")

	// 2. %-5level (左对齐5位)
	level := fmt.Sprintf("%-5s", strings.ToUpper(r.Level.String()))

	// 3. ${PID:- } (进程号，无则空)
	pid := strconv.Itoa(os.Getpid())

	// 4. [%X{traceId}/%X{spanId}] (从 Context 提取)
	var mdc []string
	if ctx != nil {
		if tid := ctx.Value("traceId"); tid != nil {
			if tidStr, ok := tid.(string); ok {
				mdc = append(mdc, tidStr)
			}
		}
		if sid := ctx.Value("spanId"); sid != nil {
			if sidStr, ok := sid.(string); ok {
				mdc = append(mdc, sidStr)
			}
		}
	}
	mdcStr := fmt.Sprintf("[%s]", strings.Join(mdc, "/"))

	// 5. [%thread] (Goroutine ID 模拟线程名)
	threadStr := fmt.Sprintf("[goroutine-%s]", getGoroutineID())

	// 6. %logger{36} (取源码包名.类名/函数名，截断36字符)
	loggerName := source(h.callerSkip)

	// 7. %msg (日志内容)
	msg := r.Message

	// 8. 处理结构化属性 (attrs)
	var attrsStr string
	if r.NumAttrs() > 0 {
		var attrs []string
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
			return true
		})
		attrsStr = " { " + strings.Join(attrs, ", ") + " }"
	}

	// 9. 拼接最终日志行
	line := fmt.Sprintf("%s %s %s --- %s %s %s : %s%s\n", timeStr, level, pid, mdcStr, threadStr, loggerName, msg, attrsStr)

	_, err := os.Stdout.WriteString(line)
	return err
}

// WithAttrs 必选继承方法
func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		Handler:    h.Handler.WithAttrs(attrs),
		writer:     h.writer,
		callerSkip: h.callerSkip,
		addSource:  h.addSource,
	}
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		Handler:    h.Handler.WithGroup(name),
		writer:     h.writer,
		callerSkip: h.callerSkip,
		addSource:  h.addSource,
	}
}

// 获取 Goroutine ID (模拟线程名)
func getGoroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := strings.Fields(string(buf[:n]))
	if len(fields) >= 2 {
		return fields[1]
	}
	return "0"
}

// 获取源码位置，模拟 %logger
func source(skip int) string {
	// 使用 runtime.Callers 手动获取调用栈
	pcs := make([]uintptr, 10)
	n := runtime.Callers(5+skip, pcs)
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return frame.Function
}
