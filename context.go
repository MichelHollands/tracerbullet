package slogTracer

import (
	"context"
)

type contextKey string

const traceHeader = contextKey("X-SlogTracer")

const setHeader = contextKey("X-SlogTracer-Set")
const defaultValue = 1

func extractSetHeader(ctx context.Context) (int, bool) {
	value, ok := ctx.Value(setHeader).(int)
	return value, ok
}

func AddSetHeaderToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, setHeader, defaultValue)
}

func AddSetHeaderWithCustomValueToContext(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, setHeader, value)
}

func extractHeader(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(traceHeader).(string)
	return value, ok
}

func AddTraceHeaderToContext(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, traceHeader, value)
}
