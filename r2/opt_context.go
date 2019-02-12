package r2

import (
	"context"
)

// OptContext sets the request context.
func OptContext(ctx context.Context) Option {
	return func(r *Request) error {
		r.Request = r.Request.WithContext(ctx)
		return nil
	}
}
