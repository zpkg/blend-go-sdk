package r2

import (
	"net/http"
	"time"
)

// Timeout sets the client timeout.
func Timeout(d time.Duration) Option {
	return func(r *Request) {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.Timeout = d
	}
}
