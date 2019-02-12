package r2

import "net/http"

// OptOnRequest sets an on request listener.
func OptOnRequest(listener func(*http.Request)) Option {
	return func(r *Request) error {
		r.OnRequest = listener
		return nil
	}
}
