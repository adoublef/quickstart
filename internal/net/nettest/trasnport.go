package nettest

import (
	"net/http"
	"net/url"
)

// Transport attempts to resolve a relative url if provided.
type Transport struct {
	http.RoundTripper
	URL string
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// if relative then use a resolve
	r2 := new(http.Request)
	*r2 = *r
	if r.URL.RawPath != "http" {
		u, err := url.Parse(t.URL)
		if err != nil {
			return nil, err
		}
		r2.URL = u.ResolveReference(r.URL)
	}
	if t.RoundTripper == nil {
		t.RoundTripper = http.DefaultTransport
	}
	return t.RoundTripper.RoundTrip(r2)
}

func WithTransport(c *http.Client, url string) *http.Client {
	c2 := new(http.Client)
	*c2 = *c
	c2.Transport = &Transport{c.Transport, url}
	return c2
}
