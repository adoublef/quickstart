package http_test

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	. "go.tmp/quickstart/internal/net/http"
	"go.tmp/quickstart/internal/testing/is"
)

func Test_AcceptHandler(t *testing.T) {
	t.Run("ContentTyp", func(t *testing.T) {
		c := newAcceptClient(t)

		type testcase struct {
			accept     string
			contentTyp string
		}

		for _, tc := range []testcase{
			{accept: "*/*", contentTyp: ContentTypJSON},
			{accept: "text/html", contentTyp: ContentTypHTML},
			{accept: "application/json", contentTyp: ContentTypJSON},
		} {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			is.OK(t, err)
			req.Header.Set("Accept", tc.accept)

			res, err := c.Do(req)
			is.OK(t, err)
			is.Equal(t, res.StatusCode, http.StatusOK)
			is.Equal(t, res.Header.Get("Content-Type"), tc.contentTyp)
		}
	})

	t.Run("ContentEncode", func(t *testing.T) {
		c := newAcceptClient(t)

		type testcase struct {
			accept string
			encode string
		}

		for _, tc := range []testcase{
			{accept: "gzip", encode: "gzip"},
			{accept: "", encode: ""},
			{accept: "identity", encode: ""},
		} {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			is.OK(t, err)
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Encoding", tc.accept)

			res, err := c.Do(req)
			is.OK(t, err)
			is.Equal(t, res.StatusCode, http.StatusOK)
			is.Equal(t, res.Header.Get("Content-Encoding"), tc.encode)
		}
	})

	t.Run("Gzip", func(t *testing.T) {
		c := newAcceptClient(t)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		is.OK(t, err)
		req.Header.Set("Accept", "text/html")
		req.Header.Set("Accept-Encoding", "gzip")

		res, err := c.Do(req)
		is.OK(t, err)

		gr, err := gzip.NewReader(res.Body)
		is.OK(t, err)

		p, err := io.ReadAll(gr) // ��)�+I�(��(�ͱ�/����+d
		is.OK(t, err)
		is.OK(t, res.Body.Close())

		is.Equal(t, string(p), "<p>text/html</p>")
	})

	t.Run("ErrNotSet", func(t *testing.T) {
		c := newAcceptClient(t)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		is.OK(t, err)

		res, err := c.Do(req)
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusNotAcceptable)
	})

	t.Run("ErrNotValid", func(t *testing.T) {
		c := newAcceptClient(t)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		is.OK(t, err)
		req.Header.Set("Accept", "application/xml")

		res, err := c.Do(req)
		is.OK(t, err)
		is.Equal(t, res.StatusCode, http.StatusNotAcceptable)
	})
}

func newAcceptClient(tb testing.TB) *TestClient {
	tb.Helper()
	// encode json data as a response
	handleTest := func() http.HandlerFunc {
		type body struct {
			Typ string `json:"contentType"`
		}

		return func(w http.ResponseWriter, r *http.Request) {
			// I could set this in the handler
			accept, ok := r.Context().Value(ContentTypOfferKey).(string)
			if !ok {
				// log error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tb.Logf("%q, ok := r.Context().Value(AcceptKey).(string)", accept)

			w.Header().Set("Content-Type", accept)
			if accept == ContentTypJSON {
				err := json.NewEncoder(w).Encode(&body{accept})
				tb.Logf("%v := json.NewEncoder(w).Encode(_)", err)
				return
			}
			// return html
			_, err := fmt.Fprintf(w, "<p>%s</p>", accept)
			tb.Logf("%v := fmt.Fprintf(w, _, accept)", err)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleTest())
	return newTestClient(tb, AcceptHandler(mux))
}
