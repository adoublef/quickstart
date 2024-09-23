package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/net/nettest"
	"go.tmp/quickstart/internal/testing/is"
)

func Test_LimitHandler(t *testing.T) {
	// for i in {1..6}; do curl http://localhost:8080/ping; done
	t.Run("OK", func(t *testing.T) {
		tc := newLimitClient(1, time.Millisecond*500)

		res, err := tc.Get("/ping")
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusNoContent)

		res, err = tc.Get("/ping")
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusTooManyRequests)

		<-time.After(time.Millisecond * 1000) // 1000ms > 500ms (ttl)

		res, err = tc.Get("/ping")
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusNoContent)
	})
}

func newLimitClient(burst int, ttl time.Duration) *http.Client {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := LimitHandler(mux, burst, ttl)
	ts := httptest.NewServer(h)

	tc := nettest.WithTransport(ts.Client(), ts.URL)
	return tc
}
