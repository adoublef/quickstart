package http

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.tmp/quickstart/internal/runtime/debug"
	// _ "golang.org/x/crypto/x509roots/fallback"
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
	h = LimitHandler(h, burst, ttl) // todo: throttle?
	h = otelhttp.NewHandler(h, "QuickStart")
	return h
}

// handleReady serves as an endpoint to signal if the service alive.
func handleReady() http.HandlerFunc {
	var (
		once sync.Once
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		// make a https request
		once.Do(func() {
			var c = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: false,
					},
				},
			}
			_, err = c.Get("https://google.com")
			debug.Printf("handleReady: _, %v = http.Get(_)", err)
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
	}
}
