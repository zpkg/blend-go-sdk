package r2

// OptBasicAuth is an option that sets the http basic auth.
func OptBasicAuth(username, password string) Option {
	return func(r *Request) error {
		r.SetBasicAuth(username, password)
		return nil
	}
}
