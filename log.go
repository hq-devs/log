package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var globalLogger *slog.Logger

// SetLogger 初始化全局日志器
func SetLogger(dir, file, unit, level string, count, size, compressT int64) error {
	filename := filepath.Join(dir, file)

	_ = os.MkdirAll(filepath.Dir(filename), 0755)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level:     getLogLevel(level),
		AddSource: true,
		// 新增：修改属性的钩子
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, ok := a.Value.Any().(*slog.Source)
				if ok {
					// 获取文件名及其上一级目录名
					// 例如: D:/.../log/log_test.go -> log/log_test.go
					dirName := filepath.Base(filepath.Dir(source.File))
					fileName := filepath.Base(source.File)
					shortPath := fmt.Sprintf("%s/%s:%d", dirName, fileName, source.Line)
					return slog.String(slog.SourceKey, shortPath)
				}
			}
			return a
		},
	})

	globalLogger = slog.New(handler)
	slog.SetDefault(globalLogger)
	return nil
}

// 核心逻辑：封装一个通用的调用函数
func logWithCaller(level slog.Level, format string, a ...interface{}) {
	if globalLogger == nil || !globalLogger.Enabled(context.Background(), level) {
		return
	}

	var msg string
	if format == "" {
		msg = fmt.Sprintln(a...)
	} else {
		msg = fmt.Sprintf(format, a...)
	}

	// 获取调用者的 PC
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // 3 是根据你的调用层级定的

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	_ = globalLogger.Handler().Handle(context.Background(), record)
}

func Debugf(format string, a ...interface{}) { logWithCaller(slog.LevelDebug, format, a...) }
func Debug(msg string, v ...interface{}) {
	logWithCaller(slog.LevelDebug, "", append([]interface{}{msg}, v...)...)
}

func Infof(format string, a ...interface{}) { logWithCaller(slog.LevelInfo, format, a...) }
func Info(msg string, v ...interface{}) {
	logWithCaller(slog.LevelInfo, "", append([]interface{}{msg}, v...)...)
}

func Warnf(format string, a ...interface{}) { logWithCaller(slog.LevelWarn, format, a...) }
func Warn(msg string, v ...interface{}) {
	logWithCaller(slog.LevelWarn, "", append([]interface{}{msg}, v...)...)
}

func Errorf(format string, a ...interface{}) { logWithCaller(slog.LevelError, format, a...) }
func Error(msg string, v ...interface{}) {
	logWithCaller(slog.LevelError, "", append([]interface{}{msg}, v...)...)
}

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
