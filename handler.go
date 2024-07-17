package slogTracer

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	handler slog.Handler
	value   int
}

func NewHandler(handler slog.Handler) Handler {
	return Handler{
		handler: handler,
		value:   defaultValue,
	}
}

func (h Handler) valid(ctx context.Context) bool {
	value, ok := extractSetHeader(ctx)
	if !ok {
		return false
	}
	return value == h.value
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.valid(ctx) && level >= slog.LevelDebug {
		return true
	}
	return h.handler.Enabled(ctx, level)
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx, record)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{
		handler: h.handler.WithAttrs(attrs),
		value:   h.value,
	}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
