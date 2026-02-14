package logger

import (
	"gonet/pkg/logger/level"
	"sync"
)

var (
	globalMu sync.RWMutex
	global   Logger
)

func ReplaceGlobals(g Logger) Logger {
	globalMu.Lock()
	defer globalMu.Unlock()
	global = g
	return g
}
func Info(msg string, fields ...Field) {
	global.Info(msg, fields...)
}
func Record() {

}
func Error(msg string, fields ...Field) {
	global.Info(msg, fields...)
}
func ParseLevel(text string) (level.Level, error) {
	return level.ParseLevel(text)
}
