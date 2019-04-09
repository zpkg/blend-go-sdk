package r2

import (
	"net/url"

	"github.com/blend/go-sdk/exception"
)

// OptURL sets the url of a request.
func OptURL(rawURL string) Option {
	return func(r *Request) error {
		var err error
		r.URL, err = url.Parse(rawURL)
		return exception.New(err)
	}
}
