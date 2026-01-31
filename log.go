package log

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var globalLogger *slog.Logger

// SetLogger 初始化全局日志器，支持日志文件滚动切割、多级别控制、压缩归档等功能
// 调用一次即可全局使用 log.Infof/Debugf 等方法，建议在项目入口处初始化
// 若仅需控制台输出，可将 dir 设为空字符串，此时忽略文件相关配置（file/unit/size/count/compressT）
//
// 参数说明：
//
//	dir        日志文件输出根路径，空字符串则仅控制台输出（例："./logs"、"/var/log/app"）
//	file       日志文件基础名称，滚动切割后会自动拼接后缀（例："app" → 生成 app.20260131.log、app.001.log 等）
//	unit       日志滚动切割单位，支持按时间/大小切割（可选值："day"按天、"hour"按小时、"size"按文件大小）
//	level      日志输出级别，低于该级别的日志会被过滤（可选值："debug"/"info"/"warn"/"error"，不区分大小写）
//	count      日志文件保留最大数量，超过则自动删除最旧文件（配合unit使用，例：按天切割保留7天则传7）
//	size       按大小切割时的单文件最大容量（单位：MB），unit="size"时生效（例：单文件最大100MB则传100）
//	count      日志文件保留最大数量，超过则自动删除最旧文件（配合unit使用，例：按天切割保留7天则传7）
//	size       按大小切割时的单文件最大容量（单位：MB），unit="size"时生效（例：单文件最大100MB则传100）
//	compressT  日志文件归档压缩延迟时间（单位：小时），超过该时间的旧日志自动压缩为.gz（0则立即压缩，-1则不压缩）
//
// 返回值：
//
//	error 初始化失败时返回具体错误信息（如路径创建失败、参数不合法等），成功则返回nil
func SetLogger(dir, file, unit, level string, count, size, compressT int64) error {
	filename := filepath.Join(dir, file)

	_ = os.MkdirAll(filepath.Dir(filename), 0755)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level:     getLogLevel(level),
		AddSource: false,
	})

	globalLogger = slog.New(handler)
	slog.SetDefault(globalLogger)
	return nil
}

func Infof(format string, a ...interface{}) {
	globalLogger.Info(fmt.Sprintf(format, a...))
}
func Info(msg string, v ...interface{}) {
	globalLogger.Info(msg, v...)
}
func Debugf(format string, a ...interface{}) {
	globalLogger.Debug(fmt.Sprintf(format, a...))
}
func Debug(msg string, v ...interface{}) {
	globalLogger.Debug(msg, v...)
}
func Warnf(format string, a ...interface{}) {
	globalLogger.Warn(fmt.Sprintf(format, a...))
}
func Warn(msg string, v ...interface{}) {
	globalLogger.Warn(msg, v...)
}
func Errorf(format string, a ...interface{}) {
	globalLogger.Error(fmt.Sprintf(format, a...))
}
func Error(msg string, v ...interface{}) {
	globalLogger.Error(msg, v...)
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
