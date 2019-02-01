package r2

import (
	"net/http"
	"time"
)

// WithTLSHandshakeTimeout sets the client transport TLSHandshakeTimeout.
func WithTLSHandshakeTimeout(d time.Duration) Option {
	return func(r *Request) {
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{
				TLSHandshakeTimeout: d,
			}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.TLSHandshakeTimeout = d
		}
	}
}
