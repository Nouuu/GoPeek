package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

var defaultLogger *Logger

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	defaultLogger = &Logger{slog.New(handler)}
}

func Default() *Logger {
	return defaultLogger
}

func (l *Logger) WithLevel(level slog.Level) *Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return &Logger{slog.New(handler)}
}
