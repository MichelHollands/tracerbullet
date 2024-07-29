package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/MichelHollands/slogTracer"
)

const header = "X-TracerBullet"
const headerValue = "abcdef"

func NewServer(logger *slog.Logger, nextHop string) http.Handler {
	mux := http.NewServeMux()

	ac := slogTracer.NewStaticAccessChecker(header, headerValue)

	client := &http.Client{
		Transport: slogTracer.NewRoundTripper(http.DefaultTransport, logger, ac),
	}

	mux.Handle("/action", slogTracer.NewMiddleware(doAction(client, logger, nextHop), logger, ac))

	var handler http.Handler = mux
	return handler
}

func doAction(client *http.Client, logger *slog.Logger, nextHop string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logger.DebugContext(ctx, "debug started")
			logger.InfoContext(ctx, "request received")

			if len(nextHop) != 0 {
				logger.DebugContext(ctx, "calling nextHop")
				r, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/action", nextHop), nil)
				logger.DebugContext(ctx, "called nextHop")
				if err != nil {
					logger.ErrorContext(ctx, fmt.Sprintf("error creating request: %v", err.Error()))
				}
				resp, err := client.Do(r)
				if err != nil {
					logger.ErrorContext(ctx, fmt.Sprintf("error calling request: %v", err.Error()))
				}
				defer resp.Body.Close()
			}

			logger.DebugContext(ctx, "debug stopped")
		},
	)
}

func main() {
	port := os.Getenv("PORT")
	nextHop := os.Getenv("NEXTHOP")

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	tracerHandler := slogTracer.NewHandler(handler)
	logger := slog.New(tracerHandler)

	logger.Info("starting server")

	srv := NewServer(logger, nextHop)
	err := http.ListenAndServe(":"+port, srv)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
