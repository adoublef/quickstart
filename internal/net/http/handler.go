package http

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Server = http.Server

var ErrServerClosed = http.ErrServerClosed

func Handler(burst int, ttl time.Duration) http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(pattern string, h http.Handler) {
		// todo: set the limit
		h = otelhttp.WithRouteTag(pattern, h)
		mux.Handle(pattern, h)
	}

	handleFunc("GET /ready", handleReady())

	// todo: csrf
	h := AcceptHandler(mux)
	h = LimitHandler(mux, burst, ttl) // todo: throttle?
	h = otelhttp.NewHandler(h, "/")
	return h
}

// handleReady serves as an endpoint to signal if the service alive.
func handleReady() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}
