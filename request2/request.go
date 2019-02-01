package request2

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/exception"
)

// MustNew returns a new request and panics on error.
func MustNew(remoteURL string, options ...Option) *Request {
	req, err := New(remoteURL, options...)
	if err != nil {
		panic(err)
	}
	return req
}

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
