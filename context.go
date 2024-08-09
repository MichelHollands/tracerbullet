package tracerbullet

import (
	"context"
)

type contextKey string

const setHeader = contextKey("X-TracerBullet-Set")
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

func extractHeader(ctx context.Context, header string) (string, bool) {
	hdr := contextKey(header)
	value, ok := ctx.Value(hdr).(string)
	return value, ok
}

func AddTraceHeaderToContext(ctx context.Context, header, value string) context.Context {
	traceHeader := contextKey(header)
	return context.WithValue(ctx, traceHeader, value)
}
