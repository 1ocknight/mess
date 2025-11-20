package logger

import (
	"context"
	"io"
	"log/slog"
)

type Logger struct {
	slog  *slog.Logger
	parse func(ctx context.Context) map[string]any
}

func New(w io.Writer, parseFunc func(ctx context.Context) map[string]any) *Logger {
	handler := slog.NewTextHandler(w, nil)
	return &Logger{
		slog:  slog.New(handler),
		parse: parseFunc,
	}
}

func (l *Logger) InfoContext(ctx context.Context, msg string) {
	attrs := []slog.Attr{}
	for k, v := range l.parse(ctx) {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.slog.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string) {
	attrs := []slog.Attr{}
	for k, v := range l.parse(ctx) {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.slog.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}
