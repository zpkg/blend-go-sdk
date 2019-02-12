package r2

import (
	"net/http"
	"time"
)

// OptOnResponse sets an on response listener.
func OptOnResponse(listener func(*http.Request, *http.Response, time.Time, error)) Option {
	return func(r *Request) error {
		r.OnResponse = listener
		return nil
	}
}
