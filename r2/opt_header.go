package r2

import (
	"net/http"
)

// Headers sets the request headers.
func Headers(headers http.Header) Option {
	return func(r *Request) {
		r.Header = headers
	}
}

// HeaderValue adds or sets a header value.
func HeaderValue(key, value string) Option {
	return func(r *Request) {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.Header.Set(key, value)
	}
}
