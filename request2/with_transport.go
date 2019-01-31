package request2

import "net/http"

// WithTransport sets the client transport for a request.
func WithTransport(transport http.RoundTripper) Option {
	return func(r *Request) {
		r.Client.Transport = transport
	}
}
