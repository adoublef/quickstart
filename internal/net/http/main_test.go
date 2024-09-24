package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/net/nettest"
)

type TestClient struct {
	*http.Client
	*nettest.Proxy
}

func newTestClient(tb testing.TB, h http.Handler) *TestClient {
	tb.Helper()

	ts := httptest.NewUnstartedServer(h)
	// https://developers.cloudflare.com/workers/platform/limits/
	// https://developers.cloudflare.com/fundamentals/reference/connection-limits/
	// https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
	// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	// subject to change
	// note: the client panics if readTimeout is less than the test timeout
	// is this a non-issue?
	ts.Config.ReadTimeout = DefaultReadTimeout
	ts.Config.WriteTimeout = DefaultWriteTimeout
	ts.Config.IdleTimeout = DefaultIdleTimeout
	ts.StartTLS()
	proxy := nettest.NewProxy(tb.Name(), strings.TrimPrefix(ts.URL, "https://"))
	tc := nettest.WithTransport(ts.Client(), "https://"+proxy.Listen())
	return &TestClient{tc, proxy}
}
