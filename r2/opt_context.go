package r2

import (
	"context"
)

// Context sets the request context.
func Context(ctx context.Context) Option {
	return func(r *Request) {
		r.Request = r.Request.WithContext(ctx)
	}
}
