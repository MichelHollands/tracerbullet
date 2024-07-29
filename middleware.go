package slogTracer

import (
	"log/slog"
	"net/http"
	"strings"
)

type middleware struct {
	handler http.Handler
	access  AccessChecker
	logger  *slog.Logger
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headerValue := r.Header.Get(m.access.Header())
	if len(headerValue) != 0 && strings.EqualFold(m.access.Value(), headerValue) {
		r = r.WithContext(AddSetHeaderToContext(r.Context()))
	}

	m.handler.ServeHTTP(w, r)
}

// NewMiddleware
func NewMiddleware(handler http.Handler, logger *slog.Logger, access AccessChecker) http.Handler {
	return &middleware{
		handler: handler,
		logger:  logger,
		access:  access,
	}
}

func NewRoundTripper(roundTripper http.RoundTripper, access AccessChecker) http.RoundTripper {
	return &addHeaderRoundTripper{
		tr:     roundTripper,
		access: access,
	}
}

type addHeaderRoundTripper struct {
	tr     http.RoundTripper
	access AccessChecker
}

func (hrt *addHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req
	ctx := req.Context()
	headerValue, ok := extractHeader(ctx)
	if ok && strings.EqualFold(hrt.access.Value(), headerValue) {
		req2 = req.Clone(ctx)
		req2.Header.Set(hrt.access.Header(), hrt.access.Value())
	}

	return hrt.tr.RoundTrip(req2)
}
