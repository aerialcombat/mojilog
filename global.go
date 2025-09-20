package mojilog

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

// InitGlobal initializes the global logger with emoji support
// This should be called once at application startup
func InitGlobal(level slog.Level, format string, addSource bool) {
	once.Do(func() {
		// Choose logger format
		switch format {
		case "json":
			// Regular JSON for machine processing
			globalLogger = SetupLogger(os.Stdout, level, format, addSource)
		case "pretty-json":
			// Pretty formatted JSON with colors
			globalLogger = SetupPrettyJSONLogger(os.Stdout, level, addSource)
		default:
			// Pretty text format (default)
			globalLogger = SetupPrettyLogger(os.Stdout, level, addSource)
		}

		// Also set as default slog logger
		slog.SetDefault(globalLogger)
	})
}

// Get returns the global logger instance
// If not initialized, it creates a default one
func Get() *slog.Logger {
	if globalLogger == nil {
		InitGlobal(slog.LevelInfo, "text", false)
	}
	return globalLogger
}

// With returns a new logger with the given attributes
func With(args ...any) *slog.Logger {
	return Get().With(args...)
}

// WithGroup returns a new logger with the given group
func WithGroup(name string) *slog.Logger {
	return Get().WithGroup(name)
}

// Debug logs at debug level
func Debug(msg string, args ...any) {
	logWithCaller(slog.LevelDebug, msg, args...)
}

// Info logs at info level
func Info(msg string, args ...any) {
	logWithCaller(slog.LevelInfo, msg, args...)
}

// Warn logs at warn level
func Warn(msg string, args ...any) {
	logWithCaller(slog.LevelWarn, msg, args...)
}

// Error logs at error level
func Error(msg string, args ...any) {
	logWithCaller(slog.LevelError, msg, args...)
}

// logWithCaller logs with the correct caller information
func logWithCaller(level slog.Level, msg string, args ...any) {
	ctx := context.TODO()
	if !Get().Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// Skip 2 frames to get the real caller:
	// 1. this function (logWithCaller)
	// 2. the wrapper function (Debug, Info, Warn, Error)
	runtime.Callers(3, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	_ = Get().Handler().Handle(ctx, r)
}

// ParseLevel converts a string to slog.Level
func ParseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Attribute convenience functions for structured logging
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

func Duration(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}