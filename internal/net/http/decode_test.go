package http_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	. "go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/net/nettest"
	"go.tmp/quickstart/internal/testing/is"
)

func Test_Decode(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		tc := newDecodeClient(t, 0, 0)

		res, err := tc.Post("/", "application/json", strings.NewReader(`{"username":"username","password":"password"}`))
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("ErrRequestEntityTooLarge", func(t *testing.T) {
		tc := newDecodeClient(t, 1, 0)

		res, err := tc.Post("/", "application/json", strings.NewReader(`{"username":"username","password":"password"}`))
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusRequestEntityTooLarge)
	})

	t.Run("ErrRequestTimeout", func(t *testing.T) {
		tc := newDecodeClient(t, 0, 1)

		err := tc.AddToxic("bandwidth", "upstream", &nettest.BandwidthToxic{Rate: 1})
		is.OK(t, err)

		s := fmt.Sprintf(`{"username":%q,"password":"password"}`, strings.Repeat("username", 1<<10))
		res, err := tc.Post("/", "application/json", strings.NewReader(s))
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusRequestTimeout)
	})

	t.Run("ErrBadRequest", func(t *testing.T) {
		c := newDecodeClient(t, 0, 0)

		type testcase struct {
			body string
			text string
		}

		for name, tc := range map[string]testcase{
			"Syntax": {
				body: `{"username:"user"}`,
				text: "invalid character 'u' at position 13",
			},
			"Syntax2": {
				body: `<"username:"user"}`,
				text: "invalid character '<' at position 1",
			},
			"Unmarshal": {
				body: `{"username":1,"password":"pass"}`,
				text: `unexpected number for field "username" at position 13`,
			},
			"Unmarshal2": {
				body: `"username:"user"}`,
				text: "unexpected string for field \"\" at position 11",
			},
			"UnknownField": {
				body: `{"never":"user"}`,
				text: `unknown field "never"`,
			},
			"Stream": {
				body: `{"username":"username","password":"password"}{}`,
				text: "request body contains more than a single JSON object",
			},
		} {
			t.Run(name, func(t *testing.T) {
				r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
				is.OK(t, err)
				r.Header.Set("Content-Type", "application/json")

				res, err := c.Do(r)
				is.OK(t, err)
				is.Equal(t, res.StatusCode, http.StatusBadRequest)
				p, err := io.ReadAll(res.Body)
				is.OK(t, err)
				is.OK(t, res.Body.Close())
				// accounting for the new-line character is stupid
				is.Equal(t, string(p), tc.text+"\n")
			})
		}
	})
}

func newDecodeClient(tb testing.TB, sz int, d time.Duration) *TestClient {
	handleTest := func() http.HandlerFunc {
		type payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := Decode[payload](w, r, sz, d)
			if err != nil {
				tb.Logf("_, %v = Decode(w, r, sz, d)", err)
				Error(w, r, err)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", handleTest())

	return newTestClient(tb, mux)
}

type Client struct {
	*http.Client
	*nettest.Proxy
}

func logf(tb testing.TB, format string, v ...any) {
	tb.Helper()
	tb.Logf(format, v...)
}
