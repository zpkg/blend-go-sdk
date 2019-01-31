package request2

import (
	"net/http"
	"time"
)

// WithResponseHeaderTimeout sets the client transport ResponseHeaderTimeout.
func WithResponseHeaderTimeout(d time.Duration) Option {
	return func(r *Request) {
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{
				ResponseHeaderTimeout: d,
			}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.ResponseHeaderTimeout = d
		}
	}
}
