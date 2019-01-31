package request2

import (
	"context"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/exception"
)

// New returns a new request.
func New(remoteURL string, options ...Option) (*Request, error) {
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, exception.New(err)
	}
	req := Request{
		Request: &http.Request{
			URL: parsedURL,
		},
	}
	for _, option := range options {
		option(&req)
	}
	return &req, nil
}

// Request is a combination of the http.Request options and the underlying client.
type Request struct {
	*http.Request
	Client *http.Client
}

// Do sends the request.
func (r *Request) Do() (*http.Response, error) {
	if r.Client != nil {
		return r.Client.Do(r.Request)
	}
	return http.DefaultClient.Do(r.Request)
}

// DoContext sends the request.
func (r *Request) DoContext(ctx context.Context) (*http.Response, error) {
	WithContext(ctx)(r)
	return r.Do()
}
