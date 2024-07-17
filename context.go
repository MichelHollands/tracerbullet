package slogTracer

import "context"

type contextKey string

const setHeader = contextKey("X-SlogTracer-Set")
const defaultValue = 1

func extractSetHeader(ctx context.Context) (int, bool) {
	value, ok := ctx.Value(setHeader).(int)
	return value, ok
}

func AddSetHeaderToContext(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, setHeader, value)
}
