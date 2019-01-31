package request2

import (
	"context"
)

// WithContext sets the request context.
func WithContext(ctx context.Context) Option {
	return func(r *Request) {
		r.Request = r.Request.WithContext(ctx)
	}
}
