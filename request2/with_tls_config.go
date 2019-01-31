package request2

import (
	"crypto/tls"
	"net/http"
)

// WithTLSConfig adds or sets a header.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(r *Request) {
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{
				TLSClientConfig: cfg,
			}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.TLSClientConfig = cfg
		}
	}
}
