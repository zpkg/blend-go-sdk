package r2

// Method sets the request method.
func Method(method string) Option {
	return func(r *Request) {
		r.Method = method
	}
}

// Get sets the request method.
func Get() Option {
	return func(r *Request) {
		r.Method = "GET"
	}
}

// Post sets the request method.
func Post() Option {
	return func(r *Request) {
		r.Method = "POST"
	}
}

// Put sets the request method.
func Put() Option {
	return func(r *Request) {
		r.Method = "PUT"
	}
}

// Delete sets the request method.
func Delete() Option {
	return func(r *Request) {
		r.Method = "DELETE"
	}
}
