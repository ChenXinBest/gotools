package log

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *slog.Logger
	ErrorLogger *slog.Logger
)

// Init 初始化日志系统
func Init() error {
	// 创建日志目录
	logDir := filepath.Join(filepath.Dir(os.Args[0]), "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 创建 info 日志文件
	infoLogFile, err := os.OpenFile(
		filepath.Join(logDir, "info.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	// 创建 error 日志文件
	errorLogFile, err := os.OpenFile(
		filepath.Join(logDir, "error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	// 配置 info 日志
	infoHandler := slog.NewTextHandler(io.MultiWriter(os.Stdout, infoLogFile), &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	})
	InfoLogger = slog.New(infoHandler)

	// 配置 error 日志
	errorHandler := slog.NewTextHandler(io.MultiWriter(os.Stderr, errorLogFile), &slog.HandlerOptions{
		Level: slog.LevelError,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	})
	ErrorLogger = slog.New(errorHandler)

	return nil
}

// Info 记录 info 级别日志
func Info(msg string, args ...any) {
	if InfoLogger == nil {
		Init()
	}
	InfoLogger.Info(msg, args...)
}

// Error 记录 error 级别日志
func Error(msg string, args ...any) {
	if ErrorLogger == nil {
		Init()
	}
	ErrorLogger.Error(msg, args...)
}

// Warn 记录 warn 级别日志
func Warn(msg string, args ...any) {
	if InfoLogger == nil {
		Init()
	}
	InfoLogger.Warn(msg, args...)
}
