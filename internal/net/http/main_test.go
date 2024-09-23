package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.tmp/quickstart/internal/net/nettest"
)

type TestClient struct {
	*http.Client
	*nettest.Proxy
}

func newTestClient(tb testing.TB, h http.Handler) *TestClient {
	tb.Helper()

	ts := httptest.NewUnstartedServer(h)
	// subject to change
	// https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
	ts.Config.ReadTimeout = 1 * time.Second
	ts.Config.WriteTimeout = 1 * time.Second
	ts.Config.IdleTimeout = 30 * time.Second
	ts.Config.ReadHeaderTimeout = 2 * time.Second
	ts.StartTLS()
	proxy := nettest.NewProxy(tb.Name(), strings.TrimPrefix(ts.URL, "https://"))
	tc := nettest.WithTransport(ts.Client(), "https://"+proxy.Listen())
	return &TestClient{tc, proxy}
}
