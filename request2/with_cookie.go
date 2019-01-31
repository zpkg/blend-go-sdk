package request2

import "net/http"

// WithCookie adds a cookie.
func WithCookie(cookie *http.Cookie) Option {
	return func(r *Request) {
		r.AddCookie(cookie)
	}
}
