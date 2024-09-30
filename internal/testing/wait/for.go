package wait

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ForHTTP
func ForHTTP(ctx context.Context, timeout time.Duration, endpoint string, opts ...func(*http.Request)) error {
	c := &http.Client{}
	o := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		// defaults added to request
		req.Header.Set("Accept", "*/*")
		for _, o := range opts {
			o(req)
		}
		res, err := c.Do(req)
		if err != nil {
			return fmt.Errorf("failed making request: %w", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return &backoff.PermanentError{Err: fmt.Errorf("not ready")}
		}
		return nil

	}
	bo := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout))
	return backoff.Retry(o, backoff.WithContext(bo, ctx))
}
