package r2

import "net/http"

// OptBasicAuth is an option that sets the http basic auth.
func OptBasicAuth(username, password string) Option {
	return func(r *Request) error {
		if r.Header == nil {
			r.Header = http.Header{}
		}
		r.SetBasicAuth(username, password)
		return nil
	}
}
