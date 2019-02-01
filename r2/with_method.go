package r2

// WithMethod sets the request method.
func WithMethod(method string) Option {
	return func(r *Request) {
		r.Method = method
	}
}

// WithMethodGet sets the request method.
func WithMethodGet() Option {
	return func(r *Request) {
		r.Method = "GET"
	}
}

// WithMethodPost sets the request method.
func WithMethodPost() Option {
	return func(r *Request) {
		r.Method = "POST"
	}
}

// WithMethodPut sets the request method.
func WithMethodPut() Option {
	return func(r *Request) {
		r.Method = "PUT"
	}
}

// WithMethodDelete sets the request method.
func WithMethodDelete() Option {
	return func(r *Request) {
		r.Method = "DELETE"
	}
}
