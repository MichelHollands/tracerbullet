package slogTracer

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})
	tracerHandler := NewHandler(handler)

	emptyCtx := context.Background()
	assert.False(t, tracerHandler.valid(emptyCtx))

	ctx := context.WithValue(context.Background(), contextKey(setHeader), 1)
	assert.True(t, tracerHandler.valid(ctx))
}

func TestLogging(t *testing.T) {
	var sb strings.Builder
	handler := slog.NewTextHandler(&sb, &slog.HandlerOptions{Level: slog.LevelError})
	tracerHandler := NewHandler(handler)
	logger := slog.New(tracerHandler)
	validCtx := AddSetHeaderToContext(context.Background())
	invalidCtx := AddSetHeaderWithCustomValueToContext(context.Background(), 2)
	emptyCtx := context.Background()

	logger.ErrorContext(validCtx, "should print")
	assert.Contains(t, sb.String(), "should print")
	sb.Reset()

	logger.ErrorContext(invalidCtx, "should print")
	assert.Contains(t, sb.String(), "should print")
	sb.Reset()

	logger.ErrorContext(emptyCtx, "should print")
	assert.Contains(t, sb.String(), "should print")
	sb.Reset()

	logger.WarnContext(validCtx, "should print")
	assert.Contains(t, sb.String(), "should print")
	sb.Reset()

	logger.WarnContext(invalidCtx, "should not print")
	assert.NotContains(t, sb.String(), "should not print")
	sb.Reset()

	logger.WarnContext(emptyCtx, "should not print")
	assert.NotContains(t, sb.String(), "should not print")
	sb.Reset()

	logger.InfoContext(validCtx, "should print")
	assert.Contains(t, sb.String(), "should print")
	sb.Reset()

	logger.InfoContext(invalidCtx, "should not print")
	assert.NotContains(t, sb.String(), "should not print")
	sb.Reset()

	logger.InfoContext(emptyCtx, "should not print")
	assert.NotContains(t, sb.String(), "should not print")
	sb.Reset()

	logger.DebugContext(validCtx, "should print")
	assert.Contains(t, sb.String(), "should print")
	sb.Reset()

	logger.DebugContext(invalidCtx, "should not print")
	assert.NotContains(t, sb.String(), "should not print")
	sb.Reset()

	logger.DebugContext(emptyCtx, "should not print")
	assert.NotContains(t, sb.String(), "should not print")
	sb.Reset()
}

func BenchmarkHandler(b *testing.B) {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	tracerHandler := NewHandler(handler)

	validCtx := AddSetHeaderToContext(context.Background())
	invalidCtx := AddSetHeaderWithCustomValueToContext(context.Background(), 2)
	emptyCtx := context.Background()

	b.Run("valid ctx", func(b *testing.B) {
		logger := slog.New(tracerHandler)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.DebugContext(validCtx, "test", "key", "label1", "key2", "label2")
		}
	})

	b.Run("invalid ctx", func(b *testing.B) {
		logger := slog.New(tracerHandler)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.DebugContext(invalidCtx, "test", "key", "label1", "key2", "label2")
		}
	})

	b.Run("empty ctx", func(b *testing.B) {
		logger := slog.New(tracerHandler)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.DebugContext(emptyCtx, "test", "key", "label1", "key2", "label2")
		}
	})

	b.Run("default handler", func(b *testing.B) {
		logger := slog.New(handler)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.DebugContext(emptyCtx, "test", "key", "label1", "key2", "label2")
		}
	})
}
