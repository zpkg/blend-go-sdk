package proxy

import (
	"net/http"
)

const (
	schemeHTTPS = "https"
)

// NewHTTPRedirect returns a new HTTPRedirect which redirects HTTP to HTTPS
func NewHTTPRedirect() *HTTPRedirect {
	return &HTTPRedirect{}
}

// HTTPRedirect redirects HTTP to HTTPS
type HTTPRedirect struct {
	BaseServer
}

// ServeHTTP redirects HTTP to HTTPS
func (hr *HTTPRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.URL.Scheme = schemeHTTPS // upgrade to https. TODO: validation on the original value?
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}

	http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
}
