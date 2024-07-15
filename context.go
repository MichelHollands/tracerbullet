package slogTracer

import "context"

type contextKey string

func extractHeader(ctx context.Context, header string) (string, bool) {
	var headerName = contextKey(header)
	header, ok := ctx.Value(headerName).(string)
	return header, ok
}

func AddToContext(ctx context.Context, header, value string) context.Context {
	return context.WithValue(ctx, contextKey(header), value)
}
