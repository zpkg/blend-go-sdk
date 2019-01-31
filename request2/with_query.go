package request2

// WithQuery adds or sets a header.
func WithQuery(key, value string) Option {
	return func(r *Request) {
		queryValues := r.URL.Query()
		queryValues.Set(key, value)
		r.URL.RawQuery = queryValues.Encode()
	}
}
