package http

import (
	"net/http"
)

type Server = http.Server

var ErrServerClosed = http.ErrServerClosed

func Handler() http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(pattern string, h http.Handler) {
		mux.Handle(pattern, h)
	}

	handleFunc("GET /ready", handleReady())
	return mux
}

// handleReady serves as an endpoint to signal if the service alive.
func handleReady() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}
