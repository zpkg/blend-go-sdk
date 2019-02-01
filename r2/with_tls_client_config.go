package r2

import (
	"crypto/tls"
	"net/http"
)

// WithTLSClientConfig adds or sets a header.
func WithTLSClientConfig(cfg *tls.Config) Option {
	return func(r *Request) {
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.TLSClientConfig = cfg
		}
	}
}
