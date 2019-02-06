package r2

import "net/http"

// Cookie adds a cookie.
func Cookie(cookie *http.Cookie) Option {
	return func(r *Request) {
		r.AddCookie(cookie)
	}
}

// CookieValue adds a cookie with a given name and value.
func CookieValue(name, value string) Option {
	return func(r *Request) {
		r.AddCookie(&http.Cookie{Name: name, Value: value})
	}
}
