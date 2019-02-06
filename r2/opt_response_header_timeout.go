package r2

import (
	"net/http"
	"time"
)

// ResponseHeaderTimeout sets the client transport ResponseHeaderTimeout.
func ResponseHeaderTimeout(d time.Duration) Option {
	return func(r *Request) {
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.ResponseHeaderTimeout = d
		}
	}
}
