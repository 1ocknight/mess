package logger

import (
	"log/slog"
)

type Logger interface {
	Info(msg string)
	Error(err error)
	With(key string, val any) Logger
}

type loggerCtxKey struct{}

var LoggerKey = loggerCtxKey{}

type log struct {
	lg *slog.Logger
}

func New(handler slog.Handler) *log {
	lg := slog.New(handler)
	return &log{
		lg: lg,
	}
}

func (l *log) Info(msg string) {
	l.lg.Info(msg)
}

func (l *log) Error(err error) {
	l.lg.Error(err.Error())
}

func (l *log) With(key string, val any) Logger {
	return &log{
		lg: l.lg.With(slog.Any(key, val)),
	}
}
