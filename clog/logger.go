package clog

import (
	"context"
	"log/slog"
	"os"
)

func replaceAttrs(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		return slog.Attr{
			Key:   "ts",
			Value: a.Value,
		}

		// other replacements may be needed in the future
	default:
		return a
	}
}

const (
	loggerNameField = "logger"
)

func NewLogger(name string, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource:   false, // use clog.WithStacktrace() at the point of logging
		Level:       level,
		ReplaceAttr: replaceAttrs,
	}
	baseHandler := slog.NewJSONHandler(os.Stderr, opts)

	logger := slog.New(baseHandler)
	logger = logger.With(slog.String(loggerNameField, name))

	return logger
}

type loggerCtxKey struct{}

func Logger(ctx context.Context) *slog.Logger {
	val := ctx.Value(loggerCtxKey{})
	if val == nil {
		return slog.Default()
	}
	logger, ok := val.(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}
