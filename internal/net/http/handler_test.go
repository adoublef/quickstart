package http_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	. "go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/testing/is"
)

func Test_handleReady(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		c, _ := newClient(t), context.Background()

		r, err := http.NewRequest("GET", "/ready", nil)
		is.OK(t, err) // return echo request

		res, err := c.Do(r)
		is.OK(t, err) // return echo response
		is.Equal(t, res.StatusCode, http.StatusOK)
	})
}

func newClient(tb testing.TB) *TestClient {
	tb.Helper()
	return newTestClient(tb, Handler(1, 200*time.Millisecond))
}
