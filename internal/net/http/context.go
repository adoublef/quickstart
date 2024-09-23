package http

import "context"

type contextKey struct{ s string }

func (k contextKey) String() string { return "go.tmp/sandbox/internal/net/http: " + k.s }

var (
	ContentTypOfferKey = &contextKey{"accept-offer"}
)

// mustValue returns the context value else panics.
func mustValue[E any](ctx context.Context, key any) E {
	v, ok := value[E](ctx, key)
	if !ok {
		panic("context is missing")
	}
	return v
}

// value returns the context value.
func value[E any](ctx context.Context, key any) (E, bool) {
	v, ok := ctx.Value(key).(E)
	return v, ok
}
