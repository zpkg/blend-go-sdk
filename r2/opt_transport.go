package r2

import "net/http"

// Transport sets the client transport for a request.
func Transport(transport http.RoundTripper) Option {
	return func(r *Request) {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.Transport = transport
	}
}
