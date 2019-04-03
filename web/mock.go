package web

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/webutil"
)

// MockCtx returns a new mock ctx.
// It is intended to be used in testing.
func MockCtx(method, path string, options ...CtxOption) *Ctx {
	return NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest(method, path), options...)
}

// Mock sends a mock request to an app.
func Mock(app *App, req *http.Request, options ...webutil.RequestOption) (*http.Response, error) {
	var err error
	for _, option := range options {
		if err = option(req); err != nil {
			return nil, err
		}
	}

	app.Server = nil
	app.Listener = nil
	app.Config.BindAddr = DefaultMockBindAddr

	startupErrors := make(chan error)
	go func() {
		if err := app.Start(); err != nil {
			startupErrors <- err
		}
	}()
	defer app.Stop()
	select {
	case <-app.NotifyStarted():
		if req.URL == nil {
			req.URL = &url.URL{}
		}
		req.URL.Host = app.Listener.Addr().String()
		return http.DefaultClient.Do(req)
	case err := <-startupErrors:
		return nil, err
	}
}

// MockMethod sends a mock request with a given method to an app.
// You should use request options to set the body of the request if it's a post or put etc.
func MockMethod(app *App, method, path string, options ...webutil.RequestOption) (*http.Response, error) {
	req := &http.Request{
		Method: method,
	}
	req.URL = &url.URL{
		Scheme: "http",
		Path:   path,
	}
	return Mock(app, req, options...)
}

// MockGet sends a mock get request to an app.
func MockGet(app *App, path string, options ...webutil.RequestOption) (*http.Response, error) {
	req := &http.Request{
		Method: "GET",
	}
	req.URL = &url.URL{
		Scheme: "http",
		Path:   path,
	}
	return Mock(app, req, options...)
}

// MockPost sends a mock post request to an app.
func MockPost(app *App, path string, body io.ReadCloser, options ...webutil.RequestOption) (*http.Response, error) {
	req := &http.Request{
		Method: "POST",
		Body:   body,
	}
	req.URL = &url.URL{
		Scheme: "http",
		Path:   path,
	}
	return Mock(app, req)
}

// MockReadBytes reads the results of a mocked request.
func MockReadBytes(res *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// MockReadBytesWithResponse reads the response of a mocked request and returns the response.
func MockReadBytesWithResponse(res *http.Response, err error) ([]byte, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	return contents, res, nil
}

// MockReadJSON reads the results of a mocked request as json.
func MockReadJSON(res *http.Response, err error) func(interface{}) error {
	return func(ref interface{}) error {
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return json.NewDecoder(res.Body).Decode(ref)
	}
}

// MockReadJSONWithResponse reads the results of a mocked request as json and also returns the response.
func MockReadJSONWithResponse(res *http.Response, err error) func(interface{}) (*http.Response, error) {
	return func(ref interface{}) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(ref); err != nil {
			return nil, err
		}
		return res, nil
	}
}

// MockReadXML reads the results of a mocked request as xml.
func MockReadXML(res *http.Response, err error) func(interface{}) error {
	return func(ref interface{}) error {
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return xml.NewDecoder(res.Body).Decode(ref)
	}
}

// MockReadXMLWithResponse reads the results of a mocked request as xml and also returns the response.
func MockReadXMLWithResponse(res *http.Response, err error) func(interface{}) (*http.Response, error) {
	return func(ref interface{}) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if err := xml.NewDecoder(res.Body).Decode(ref); err != nil {
			return nil, err
		}
		return res, nil
	}
}

// MockReadDiscard discards the results of a mocked request.
func MockReadDiscard(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	if _, err = io.Copy(ioutil.Discard, res.Body); err != nil {
		return err
	}
	return nil
}

// MockReadDiscardWithResponse discards the results of a mocked request and returns the response.
func MockReadDiscardWithResponse(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(ioutil.Discard, res.Body); err != nil {
		return nil, err
	}
	return res, nil
}
