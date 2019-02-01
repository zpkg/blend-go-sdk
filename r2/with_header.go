package r2

import (
	"net/http"
)

// WithHeader adds or sets a header.
func WithHeader(key, value string) Option {
	return func(r *Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(key, value)
	}
}
