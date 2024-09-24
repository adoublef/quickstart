package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

var tooManyRequestsHandler = &statusError{
	http.StatusTooManyRequests,
	fmt.Sprintf(`
<html>
  <head>
    <title>Too Many Requests</title>
  </head>
  <body>
    <h1>Too Many Requests</h1>
    <p>You're doing that too often! Try again later.</p>
  </body>
</html>`[1:]),
}

func LimitHandler(h http.Handler, burst int, ttl time.Duration) http.Handler {
	var l = rate.NewLimiter(rate.Every(ttl), burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("RateLimit-Limit", strconv.FormatFloat(float64(l.Limit()), 'f', 0, 64))
		w.Header().Add("RateLimit-Reset", "1") // note: 1 is fixed
		w.Header().Add("RateLimit-Remaining", strconv.FormatFloat(l.Tokens(), 'f', 0, 64))
		// throttle safe requests and limit non-safe requests
		if !l.Allow() {
			tooManyRequestsHandler.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
