package slogTracer

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})
	tracerHandler := NewHandler(handler, "X-SlogTracer", "abc")

	emptyCtx := context.Background()
	assert.False(t, tracerHandler.valid(emptyCtx))

	ctx := context.WithValue(context.Background(), contextKey("X-SlogTracer"), "abc")
	assert.True(t, tracerHandler.valid(ctx))
}

func TestLogging(t *testing.T) {
	var sb strings.Builder
	handler := slog.NewTextHandler(&sb, &slog.HandlerOptions{Level: slog.LevelError})
	tracerHandler := NewHandler(handler, "X-SlogTracer", "abc")
	logger := slog.New(tracerHandler)
	validCtx := context.WithValue(context.Background(), contextKey("X-SlogTracer"), "abc")
	invalidCtx := context.WithValue(context.Background(), contextKey("X-SlogTracer"), "def")
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
