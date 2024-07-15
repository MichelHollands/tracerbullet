package slogTracer

import (
	"context"
	"log/slog"
	"strings"
)

type contextKey string

func extractHeader(ctx context.Context, header string) (string, bool) {
	var headerName = contextKey(header)
	header, ok := ctx.Value(headerName).(string)
	return header, ok
}

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	handler slog.Handler
	header  string
	value   string
}

func NewHandler(handler slog.Handler, headerName, expected string) Handler {
	return Handler{
		handler: handler,
		header:  headerName,
		value:   expected,
	}
}

func (h Handler) valid(ctx context.Context) bool {
	header, ok := extractHeader(ctx, h.header)
	if !ok {
		return false
	}
	// HTTP headers are case insensitive
	if strings.EqualFold(header, h.value) {
		return true
	}
	return false
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
