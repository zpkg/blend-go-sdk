package secrets

import (
	"net/http"
	"strconv"
)

// Option a thing that we can do to modify a request.
type Option func(req *http.Request)

// Version adds a version to the request.
func Version(version int) Option {
	return func(req *http.Request) {
		req.URL.Query().Add("version", strconv.Itoa(version))
	}
}
