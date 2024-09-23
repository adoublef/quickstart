package rate

import (
	"errors"
	"fmt"
	"time"

	"go.tmp/quickstart/internal/runtime/debug"
)

var (
	ErrDivByZero = errors.New("division by zero")
	ErrInvalid   = errors.New("fraction invalid")
)

type Duration = time.Duration

type Rate struct {
	N int           // burst
	D time.Duration // duration
}

func (r Rate) MarshalText() ([]byte, error) {
	// note: duration cannot be negative
	if r.D < 0 {
		return nil, fmt.Errorf("duration is negative: %w", ErrInvalid)
	}
	return []byte(r.String()), nil
}

func (r *Rate) UnmarshalText(p []byte) (err error) {
	*r, err = ParseRate(string(p))
	debug.Printf("_, %s = ParseRate(string(p))", err)
	if err == nil {
		return nil
	}
	return err
}

func (r Rate) String() string { return fmt.Sprintf("%d/%s", r.N, r.D) }

func ParseRate(s string) (Rate, error) {
	var n int
	var b string
	np, err := fmt.Sscanf(s, "%d/%s", &n, &b)
	if err != nil {
		// note: wrap
		return Rate{}, fmt.Errorf("%s: %w", err, ErrInvalid)
	}
	if np != 2 {
		return Rate{}, ErrInvalid
	}
	d, err := time.ParseDuration(b)
	if err != nil {
		return Rate{}, fmt.Errorf("%s: %w", err, ErrInvalid)
	}
	// note: cannot be negative
	if d < 0 {
		return Rate{}, fmt.Errorf("duration is negative: %w", ErrInvalid)
	}
	return Rate{n, d}, nil
}
