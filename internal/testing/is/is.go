package is

import (
	"errors"
	"testing"

	"go.adoublef.dev/is"
)

// OK fails test if error is not nil
func OK(tb testing.TB, err error) {
	is := is.NewRelaxed(tb)
	is.Helper()
	is.NoErr(err)
}

// NotOK fails test if errors don't match
func NotOK(tb testing.TB, err error, target any) {
	is := is.NewRelaxed(tb)
	is.Helper()
	switch v := target.(type) {
	case nil:
		is.True(errors.Is(err, nil))
	case error:
		is.True(errors.Is(err, v))
	default:
		is.True(errors.As(err, v))
	}
}

// Equal fails test if two values are not Equal
func Equal[V any](tb testing.TB, a, b V) {
	is := is.NewRelaxed(tb)
	is.Helper()
	is.Equal(a, b)
}

// True fails test if expression is false
func True(tb testing.TB, exp bool) {
	is := is.NewRelaxed(tb)
	is.Helper()
	is.True(exp)
}
