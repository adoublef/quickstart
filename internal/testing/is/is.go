package is

import (
	"errors"
	"testing"

	"go.adoublef.dev/is"
)

func OK(tb testing.TB, err error) {
	is := is.NewRelaxed(tb)
	is.Helper()
	is.NoErr(err)
}

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
