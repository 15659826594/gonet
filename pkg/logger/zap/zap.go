package zap

import (
	"fmt"
	"gonet/pkg/config"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func New(config *config.Log) *zap.Logger {
	// 使用 zap 的 NewDevelopmentConfig 快速配置
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	// level 默认是0 是 InfoLevel
	level, _ := zapcore.ParseLevel(config.Level)
	// 创建 Encoder
	encoder := &logEncoder{
		Encoder: zapcore.NewJSONEncoder(cfg.EncoderConfig), // 生产环境使用json格式
		config:  config,
	}
	if gin.IsDebugging() {
		encoder.Encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}
	// 创建 Logger
	return zap.New(zapcore.NewCore(
		encoder,                    // 使用自定义的 Encoder
		zapcore.AddSync(os.Stdout), // 输出到控制台
		level,                      // 设置日志级别
	), zap.AddCaller(), zap.AddCallerSkip(2))
}

// logEncoder 时间分片和level分片同时做
type logEncoder struct {
	zapcore.Encoder
	errFile     *os.File
	file        *os.File
	currentDate string
	config      *config.Log
}

func (e *logEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 先调用原始的 EncodeEntry 方法生成日志行
	buff, err := e.Encoder.EncodeEntry(entry, fields)
	if gin.IsDebugging() || e.config.Type != "File" {
		return buff, err
	}
	if err != nil {
		return nil, err
	}
	data := buff.String()
	buff.Reset()
	if e.config.Prefix != "" {
		buff.AppendString(fmt.Sprintf("[%s] ", e.config.Prefix) + data)
		data = buff.String()
	}
	// 时间分片
	now := time.Now().Format(time.DateOnly)
	if e.currentDate != now {
		os.MkdirAll(fmt.Sprintf("%s/%s", e.config.Path, now), 0666)
		// 时间不同，先创建目录
		name := fmt.Sprintf("%s/%s/out.log", e.config.Path, now)
		file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		e.file = file
		e.currentDate = now
	}

	switch entry.Level {
	case zapcore.ErrorLevel:
		if e.errFile == nil {
			name := fmt.Sprintf("%s/%s/%s.log", e.config.Path, now, entry.Level.String())
			file, _ := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			e.errFile = file
		}
		e.errFile.WriteString(data)
	default:
	}

	if e.currentDate == now {
		e.file.WriteString(data)
	}
	return buff, nil
}
