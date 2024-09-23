package http

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

const (
	StatusTooManyRequests = http.StatusTooManyRequests
)

func LimitHandler(h http.Handler, burst int, ttl time.Duration) http.Handler {
	var l = rate.NewLimiter(rate.Every(ttl), burst)
	// fixme: remove
	tooManyReq := func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(StatusTooManyRequests), StatusTooManyRequests)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("RateLimit-Limit", strconv.FormatFloat(float64(l.Limit()), 'f', 0, 64))
		w.Header().Add("RateLimit-Reset", "1") // note: 1 is fixed
		w.Header().Add("RateLimit-Remaining", strconv.FormatFloat(l.Tokens(), 'f', 0, 64))
		// throttle safe requests and limit non-safe requests
		if !l.Allow() {
			tooManyReq(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
