package r2

import "net/url"

// Query set the fully querystring.
func Query(query url.Values) Option {
	return func(r *Request) {
		r.URL.RawQuery = query.Encode()
	}
}

// QueryValue adds or sets a query value.
func QueryValue(key, value string) Option {
	return func(r *Request) {
		queryValues := r.URL.Query()
		queryValues.Set(key, value)
		r.URL.RawQuery = queryValues.Encode()
	}
}
