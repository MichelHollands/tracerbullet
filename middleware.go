package tracerbullet

import (
	"log/slog"
	"net/http"
	"strings"
)

type middleware struct {
	handler http.Handler
	logger  *slog.Logger
	access  AccessChecker
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headerValue := r.Header.Get(m.access.Header())
	if len(headerValue) != 0 && m.access.Check(headerValue) {
		// Used by handler
		r = r.WithContext(AddSetHeaderToContext(r.Context()))
		// The original header
		r = r.WithContext(AddTraceHeaderToContext(r.Context(), m.access.Header(), m.access.Value()))
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

func NewRoundTripper(roundTripper http.RoundTripper, logger *slog.Logger, access AccessChecker) http.RoundTripper {
	return &addHeaderRoundTripper{
		tr:     roundTripper,
		logger: logger,
		access: access,
	}
}

type addHeaderRoundTripper struct {
	tr     http.RoundTripper
	logger *slog.Logger
	access AccessChecker
}

func (hrt *addHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req
	ctx := req.Context()
	headerValue, ok := extractHeader(ctx, hrt.access.Header())
	if ok && strings.EqualFold(hrt.access.Value(), headerValue) {
		req2 = req.Clone(ctx)
		req2.Header.Set(hrt.access.Header(), hrt.access.Value())
	}

	return hrt.tr.RoundTrip(req2)
}
