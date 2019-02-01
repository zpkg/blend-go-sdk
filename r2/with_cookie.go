package r2

import "net/http"

// WithCookie adds a cookie.
func WithCookie(cookie *http.Cookie) Option {
	return func(r *Request) {
		r.AddCookie(cookie)
	}
}

// WithCookieValue adds a cookie with a given name and value.
func WithCookieValue(name, value string) Option {
	return func(r *Request) {
		r.AddCookie(&http.Cookie{Name: name, Value: value})
	}
}
