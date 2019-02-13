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

	// OnRequest and OnResponse are lifecycle hooks.
	OnRequest  func(*http.Request)
	OnResponse func(*http.Request, *http.Response, time.Time, error)
}

// Do executes the request.
func (r *Request) Do() (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	var err error
	started := time.Now().UTC()

	if r.OnRequest != nil {
		r.OnRequest(r.Request)
	}

	var res *http.Response
	if r.Client != nil {
		res, err = r.Client.Do(r.Request)
	} else {
		res, err = http.DefaultClient.Do(r.Request)
	}

	if r.OnResponse != nil {
		r.OnResponse(r.Request, res, started, err)
	}
	return res, err
}

// Close executes and closes the response.
// It returns the response for metadata purposes.
func (r *Request) Close() (*http.Response, error) {
	res, err := r.Do()
	if err != nil {
		return nil, err
	}
	return res, exception.New(res.Body.Close())
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
	return string(contents), err
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
