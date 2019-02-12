package r2

import (
	"net/http"
	"net/url"
	"time"
)

// New returns a new request.
// The default method is GET.
func New(remoteURL string, options ...Option) (*Request, error) {
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, err
	}

	r := &Request{
		Request: &http.Request{
			Method: MethodGet, // default to get
			URL:    parsedURL,
		},
	}

	for _, option := range options {
		if err = option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Request is a combination of the http.Request options and the underlying client.
type Request struct {
	*http.Request
	Client *http.Client

	OnRequest  func(*http.Request)
	OnResponse func(*http.Request, *http.Response, time.Time, error)
}
