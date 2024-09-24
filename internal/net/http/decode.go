package http

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"go.tmp/quickstart/internal/runtime/debug"
)

const (
	DefaultReadTimeout    = 100 * time.Second // Cloudflare's default read request timeout of 100s
	DefaultWriteTimeout   = 30 * time.Second  // Cloudflare's default write request timeout of 30s
	DefaultIdleTimeout    = 900 * time.Second // Cloudflare's default write request timeout of 900s
	DefaultMaxHeaderBytes = 32 * (1 << 10)
	DefaultMaxBytes       = 1 << 20 // Cloudflare's free tier limits of 100mb
)

func Decode[T any](w http.ResponseWriter, r *http.Request, sz int, d time.Duration) (v T, err error) {
	if r.Body == nil {
		return v, unauthorizedError("request body is missing")
	}
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")) // params could include versioning
	if err != nil || !(mt == ContentTypJSON) {
		return v, unsupportedMediaTypeError("content-type is not application/json")
	}
	if sz > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, int64(sz))
	}
	if d > 0 {
		rc := http.NewResponseController(w)
		err = rc.SetReadDeadline(time.Now().Add(d))
		if err != nil {
			// note: if action not allowed, should maybe wrap this
			// ErrNotSupported or Deadline error
			// which is caused by how?
			return v, err
		}
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // important for unknown fields
	if err := dec.Decode(&v); err != nil {
		// note: error will not be forwarded to the client
		debug.Printf("%v := dec.Decode(&v)", err)
		switch v := *new(T); {
		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.As(err, new(*json.SyntaxError)):
			se := err.(*json.SyntaxError)
			ch, _ := utf8.DecodeRune([]byte(se.Error()[19:]))
			return v, badRequestError("invalid character '%c' at position %d", ch, se.Offset)
		case errors.As(err, new(*json.UnmarshalTypeError)):
			e := err.(*json.UnmarshalTypeError)
			return v, badRequestError("unexpected %s for field %q at position %d", e.Value, e.Field, e.Offset)
		// There is an open issue at https://github.com/golang/go/issues/29035
		// regarding turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			return v, badRequestError("unknown field %s", err.Error()[20:])
		// An io.EOF error is returned by Decode() if the request body is empty.
		case errors.Is(err, io.EOF):
			return v, unauthorizedError("request body could not be read properly")
		case errors.As(err, new(*http.MaxBytesError)):
			return v, requestEntityTooLargeError("maximum allowed request size is %s", strconv.Itoa(sz))
		case errors.As(err, new(*net.OpError)):
			return v, requestTimeoutError("failed to process request in time, please try again")
		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response. May want to wrap this error.
		default:
			return v, err
		}
	}
	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	if err = dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		// fixme: 4xx
		return *new(T), badRequestError("request body contains more than a single JSON object")
	}
	return v, nil
}
