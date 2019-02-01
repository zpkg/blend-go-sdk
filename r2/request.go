package r2

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/exception"
)

// New returns a new request.
// The default method is GET.
func New(remoteURL string, options ...Option) *Request {
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return &Request{
			Err: err,
		}
	}
	req := &Request{
		Request: &http.Request{
			Method: MethodGet, // default to get
			URL:    parsedURL,
		},
	}
	return req.WithOptions(options...)
}

// Request is a combination of the http.Request options and the underlying client.
type Request struct {
	*http.Request
	Client *http.Client
	Err    error
}

// WithOptions applies a given set of options.
func (r *Request) WithOptions(options ...Option) *Request {
	for _, option := range options {
		option(r)
	}
	return r
}

// Do sends the request.
func (r *Request) Do() (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if r.Client != nil {
		return r.Client.Do(r.Request)
	}
	return http.DefaultClient.Do(r.Request)
}

// Discard discards the response of a request.
func (r *Request) Discard() error {
	res, err := r.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(ioutil.Discard, res.Body)
	return exception.New(err)
}

// CopyTo copies the response body to a given writer.
func (r *Request) CopyTo(dst io.Writer) (int64, error) {
	res, err := r.Do()
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return io.Copy(dst, res.Body)
}

// Bytes returns the contents of the response as a byte array.
func (r *Request) Bytes() ([]byte, error) {
	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// JSON reads a response body and decodes it into a given object.
func (r *Request) JSON(ref interface{}) error {
	res, err := r.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(ref); err != nil {
		return err
	}
	return nil
}

// XML reads a response body and decodes it into a given object.
func (r *Request) XML(ref interface{}) error {
	res, err := r.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := xml.NewDecoder(res.Body).Decode(ref); err != nil {
		return err
	}
	return nil
}
