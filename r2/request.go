package r2

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/exception"
)

// New returns a new request.
// The default method is GET.
func New(remoteURL string, options ...Option) *Request {
	var r Request
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		r.Err = err
		return &r
	}
	r.Request = &http.Request{
		Method: MethodGet,
		URL:    parsedURL,
	}
	for _, option := range options {
		if err = option(&r); err != nil {
			r.Err = err
			return &r
		}
	}
	return &r
}

// Request is a combination of the http.Request options and the underlying client.
type Request struct {
	*http.Request
	Client *http.Client

	// Err is an error set on construction.
	// It pre-empts the request going out.
	Err error

	// ResponseBodyInterceptor is an optional custom step to alter the response stream.
	ResponseBodyInterceptor ReaderInterceptor

	// OnRequest and OnResponse are lifecycle hooks.
	OnRequest  []OnRequestListener
	OnResponse []OnResponseListener
}

// Do executes the request.
func (r *Request) Do() (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var err error
	started := time.Now().UTC()

	for _, listener := range r.OnRequest {
		if err = listener(r.Request); err != nil {
			return nil, err
		}
	}

	var res *http.Response
	if r.Client != nil {
		res, err = r.Client.Do(r.Request)
	} else {
		res, err = http.DefaultClient.Do(r.Request)
	}
	for _, listener := range r.OnResponse {
		if err = listener(r.Request, res, started, err); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	// apply the interceptor if supplied.
	res.Body = r.responseBody(res)
	return res, nil
}

// Close executes and closes the response.
// It returns the response for metadata purposes.
// It does not read any data from the response.
func (r *Request) Close() (*http.Response, error) {
	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	return res, exception.New(res.Body.Close())
}

// Discard reads the response fully and discards all data it reads.
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
	count, err := io.Copy(dst, res.Body)
	if err != nil {
		return count, exception.New(err)
	}
	return count, nil
}

// Bytes reads the response and returns it as a byte array.
func (r *Request) Bytes() ([]byte, error) {
	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, exception.New(err)
	}
	return contents, nil
}

// String reads the response and returns it as a string
func (r *Request) String() (string, error) {
	contents, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

// JSON reads the response as json into a given object.
func (r *Request) JSON(dst interface{}) error {
	res, err := r.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return exception.New(json.NewDecoder(res.Body).Decode(dst))
}

// XML reads the response as json into a given object.
func (r *Request) XML(dst interface{}) error {
	res, err := r.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return exception.New(xml.NewDecoder(res.Body).Decode(dst))
}

//
// utils
//

// responseBody applies a ResponseBodyInterceptor if it's supplied.
func (r *Request) responseBody(res *http.Response) io.ReadCloser {
	if r.ResponseBodyInterceptor != nil {
		return NewReadCloser(res.Body, r.ResponseBodyInterceptor)
	}
	return res.Body
}
