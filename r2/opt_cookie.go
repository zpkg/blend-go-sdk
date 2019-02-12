package r2

import "net/http"

// OptCookie adds a cookie.
func OptCookie(cookie *http.Cookie) Option {
	return func(r *Request) error {
		r.AddCookie(cookie)
		return nil
	}
}

// OptCookieValue adds a cookie with a given name and value.
func OptCookieValue(name, value string) Option {
	return func(r *Request) error {
		r.AddCookie(&http.Cookie{Name: name, Value: value})
		return nil
	}
}
