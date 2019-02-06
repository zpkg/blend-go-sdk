package r2

import "io"

// Body sets the post body on the request.
func Body(contents io.ReadCloser) Option {
	return func(r *Request) {
		r.Body = contents
	}
}
