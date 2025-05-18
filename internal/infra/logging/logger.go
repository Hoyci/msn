package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type ctxKey struct{}

var fallbackLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return logger
	}
	return fallbackLogger
}

func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func NewLogger(w io.Writer, env string) *slog.Logger {
	var handler slog.Handler

	if env == "development" {
		handler = NewPrettyHandler(w, slog.LevelDebug)
	} else {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(handler)
}
