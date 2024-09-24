package http

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"

	"github.com/golang/gddo/httputil"
	"go.tmp/quickstart/internal/runtime/debug"
)

const (
	ContentTypJSON = "application/json"
	ContentTypHTML = "text/html"

	StatusNotAcceptable = http.StatusNotAcceptable
)

// AcceptHandler
func AcceptHandler(h http.Handler) http.Handler {
	var (
		ct = []string{ContentTypJSON, ContentTypHTML}
		ce = []string{"identity", "gzip" /* "deflate", "zstandard", "zlib" */}
	)
	// note: if we allow compression option
	// panic if invalid
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accept := httputil.NegotiateContentType(r, ct, "")
		if accept == "" {
			w.WriteHeader(StatusNotAcceptable) // todo: [http.Error]
			return
		}
		debug.Printf("AcceptHandler: %q = httputil.NegotiateContentType(r, ct, _)", accept)

		ctx = context.WithValue(ctx, ContentTypOfferKey, accept)
		encoding := httputil.NegotiateContentEncoding(r, ce)
		// note: should we always encode text/html?
		if encoding == "" || encoding == ce[0] {
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		debug.Printf("AcceptHandler: %q = negotiate.ContentEncoding(r, ce)", encoding)

		w.Header().Set("Content-Encoding", encoding)
		//
		gw, _ := gzip.NewWriterLevel(w, gzip.DefaultCompression)
		defer gw.Close() // should I defer?
		h.ServeHTTP(&gzipWriter{w, gw}, r.WithContext(ctx))
	})
}

type gzipWriter struct {
	http.ResponseWriter
	io.Writer
}

// Write implements http.ResponseWriter.
func (w *gzipWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}
