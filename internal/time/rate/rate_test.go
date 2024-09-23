package rate_test

import (
	"testing"

	"go.tmp/quickstart/internal/testing/is"
	. "go.tmp/quickstart/internal/time/rate"
)

func Test_ParseFraction(t *testing.T) {
	type testcase struct {
		in   string
		want error
	}

	var tt = map[string]testcase{
		"1/2": {
			in: "1/1s",
		},
		"ErrInvalidNegative": {
			in:   "1/-1s",
			want: ErrInvalid,
		},
		"0/0s": {
			in: "0/0s",
		},
		"ErrInvalidBurst": {
			in:   "a/b",
			want: ErrInvalid,
		},
		"ErrInvalidDuration": {
			in:   "1/b",
			want: ErrInvalid,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := ParseRate(tc.in)
			is.NotOK(t, err, tc.want)
			t.Logf("_, %v = ParseFraction(%s)", err, tc.in)
		})
	}
}
