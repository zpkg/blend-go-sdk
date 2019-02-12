package r2

import "net/url"

// OptQuery set the fully querystring.
func OptQuery(query url.Values) Option {
	return func(r *Request) error {
		r.URL.RawQuery = query.Encode()
		return nil
	}
}

// OptQueryValue adds or sets a query value.
func OptQueryValue(key, value string) Option {
	return func(r *Request) error {
		queryValues := r.URL.Query()
		queryValues.Set(key, value)
		r.URL.RawQuery = queryValues.Encode()
		return nil
	}
}
