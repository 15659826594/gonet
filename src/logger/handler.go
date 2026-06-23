package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gta-spec/utils/slog"
)

var (
	wd string
)

func init() {
	wd, _ = os.Getwd()
}

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var defaultLogFormatter = func(cfg Config, ctx context.Context, r slog.Record) error {
	levelColor := LevelColor(r.Level)

	var traceId string
	var spanId string
	if ctx != nil {
		if tid := ctx.Value("traceId"); tid != nil {
			traceId, _ = tid.(string)
		}
		if sid := ctx.Value("spanId"); sid != nil {
			spanId, _ = sid.(string)
		}
	}

	// 获取一层调用栈信息
	var source string
	if cfg.addSource {
		caller := Source(cfg.callerSkip)
		source = fmt.Sprintf("%s:%d", caller.File, caller.Line)
	}

	var attrs string
	if r.NumAttrs() > 0 {
		var arr []string
		r.Attrs(func(a slog.Attr) bool {
			//bytes, err := json.MarshalIndent(a.Value.Any(), "", "\t")
			//if err == nil {
			//	fmt.Println(string(bytes))
			//}
			arr = append(arr, fmt.Sprintf("%s=%v", a.Key, a.Value))
			return true
		})
		attrs = strings.Join(arr, " ")
	}

	var err error
	for _, wri := range cfg.writer {
		var level string
		time := r.Time.Format(time.DateTime + ".000")
		switch wri.(type) {
		case *_slog.Roller:
			level = fmt.Sprintf(" %-5s |", strings.ToUpper(r.Level.String()))
		default:
			level = fmt.Sprintf("%s %-5s %s|", levelColor, strings.ToUpper(r.Level.String()), reset)
			if traceId != "" {
				traceId = fmt.Sprintf("%s %-30s %s|", magenta, traceId, reset)
			}
			if spanId != "" {
				spanId = fmt.Sprintf("%s %-5s %s|", blue, spanId, reset)
			}
		}
		_, err = fmt.Fprint(wri, fmt.Sprintf("[GIN] %v |%s%s%s %s \n%s %s\n",
			time,
			level,
			traceId,
			spanId,
			source,
			r.Message,
			attrs,
		))
	}
	return err
}

// Source 返回日志事件的事件源
func Source(skip int) *slog.Source {
	var pcs [1]uintptr
	n := runtime.Callers(skip+6, pcs[:])
	if n == 0 {
		return nil
	}

	fs := runtime.CallersFrames(pcs[:n])
	f, _ := fs.Next()
	file, _ := filepath.Rel(wd, f.File)
	return &slog.Source{
		Function: f.Function,
		File:     filepath.ToSlash(file),
		Line:     f.Line,
	}
}

type Log struct {
	slog.Handler
	config Config
}

func (h *Log) Handle(ctx context.Context, r slog.Record) error {
	if h.config.formatter != nil {
		return h.config.formatter(h.config, ctx, r) // 使用自定义格式化器
	}
	return defaultLogFormatter(h.config, ctx, r)
}

func (h *Log) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Log{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *Log) WithGroup(name string) slog.Handler {
	return &Log{Handler: h.Handler.WithGroup(name)}
}

type Config struct {
	callerSkip  int
	writer      []io.Writer
	addSource   bool
	level       slog.Leveler
	replaceAttr func(groups []string, a slog.Attr) slog.Attr
	formatter   func(Config, context.Context, slog.Record) error
}

// LevelColor is the ANSI color for level
func LevelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return blue
	case slog.LevelInfo:
		return green
	case slog.LevelWarn:
		return yellow
	case slog.LevelError:
		return red
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (c *Config) ResetColor() string {
	return reset
}

type Option func(*Config)

func WithCallerSkip(skip int) Option {
	return func(c *Config) {
		c.callerSkip = skip
	}
}

func WithWriter(w ...io.Writer) Option {
	return func(c *Config) {
		c.writer = w
	}
}

func WithAddSource(add bool) Option {
	return func(c *Config) {
		c.addSource = add
	}
}

func WithLevel(level slog.Leveler) Option {
	return func(c *Config) {
		c.level = level
	}
}

func WithReplaceAttr(fn func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(c *Config) {
		c.replaceAttr = fn
	}
}

func WithFormatter(fn func(Config, context.Context, slog.Record) error) Option {
	return func(c *Config) {
		c.formatter = fn
	}
}

func New(opts ...Option) *slog.Logger {
	c := &Config{
		callerSkip:  0,
		writer:      []io.Writer{os.Stdout},
		addSource:   true,
		level:       slog.LevelInfo,
		replaceAttr: nil,
		formatter:   defaultLogFormatter,
	}

	for _, opt := range opts {
		opt(c)
	}

	return slog.New(&Log{
		Handler: slog.NewTextHandler(nil, &slog.HandlerOptions{
			AddSource:   c.addSource,
			Level:       c.level,
			ReplaceAttr: c.replaceAttr,
		}),
		config: *c,
	})
}
