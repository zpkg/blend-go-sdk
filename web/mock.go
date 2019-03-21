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
func Mock(app *App, req *http.Request, options ...func(*http.Request) error) (*http.Response, error) {
	var err error
	for _, option := range options {
		if err = option(req); err != nil {
			return nil, err
		}
	}

	app.Config.BindAddr = DefaultMockBindAddr

	if err := app.Start(); err != nil {
		return nil, err
	}
	defer app.Stop()
	req.Host = app.Listener.Addr().String()
	return http.DefaultClient.Do(req)
}

// MockGet sends a mock get request to an app.
func MockGet(app *App, path string, options ...func(*http.Request) error) (*http.Response, error) {
	req := &http.Request{
		Method: "GET",
	}
	req.URL = &url.URL{
		Scheme: "http",
		Host:   app.Listener.Addr().String(),
		Path:   path,
	}
	return Mock(app, req, options...)
}

// MockPost sends a mock post request to an app.
func MockPost(app *App, path string, body io.ReadCloser, options ...func(*http.Request) error) (*http.Response, error) {
	req := &http.Request{
		Method: "POST",
		Body:   body,
	}
	req.URL = &url.URL{
		Scheme: "http",
		Host:   app.Listener.Addr().String(),
		Path:   path,
	}
	return Mock(app, req)
}

// MockBytes reads the results of a mocked request.
func MockBytes(res *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// MockJSON reads the results of a mocked request as json.
func MockJSON(res *http.Response, err error) func(interface{}) error {
	return func(ref interface{}) error {
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return json.NewDecoder(res.Body).Decode(ref)
	}
}

// MockXML reads the results of a mocked request as xml.
func MockXML(res *http.Response, err error) func(interface{}) error {
	return func(ref interface{}) error {
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return xml.NewDecoder(res.Body).Decode(ref)
	}
}

// MockDiscard discards the results of a mocked request.
func MockDiscard(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	if _, err = io.Copy(ioutil.Discard, res.Body); err != nil {
		return err
	}
	return nil
}
