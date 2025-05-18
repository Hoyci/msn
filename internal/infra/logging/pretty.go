package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
	writer   io.Writer
	minLevel slog.Level
	attrs    []slog.Attr
}

func NewPrettyHandler(w io.Writer, minLevel slog.Level) *PrettyHandler {
	return &PrettyHandler{
		writer:   w,
		minLevel: minLevel,
	}
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.minLevel
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := h.levelString(r.Level)

	fmt.Fprintf(h.writer, "%s | %s | %s", timeStr, levelStr, r.Message)

	for _, attr := range h.attrs {
		fmt.Fprintf(h.writer, " %s=%v", color.GreenString(attr.Key), attr.Value.Any())
	}

	r.Attrs(func(attr slog.Attr) bool {
		fmt.Fprintf(h.writer, " %s=%v", color.GreenString(attr.Key), attr.Value.Any())
		return true
	})

	fmt.Fprintln(h.writer)
	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		writer:   h.writer,
		minLevel: h.minLevel,
		attrs:    append(h.attrs, attrs...),
	}
}

func (h *PrettyHandler) levelString(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return color.MagentaString("DEBUG")
	case slog.LevelInfo:
		return color.CyanString("INFO")
	case slog.LevelWarn:
		return color.YellowString("WARN")
	case slog.LevelError:
		return color.RedString("ERROR")
	default:
		return level.String()
	}
}
