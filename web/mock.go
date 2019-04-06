package web

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

// Mock sends a mock request to an app.
// It will reset the app Server, Listener, and will set the request host to the listener address
// for a randomized local listener.
func Mock(app *App, req *http.Request, options ...webutil.RequestOption) *MockResult {
	var err error
	for _, option := range options {
		if err = option(req); err != nil {
			return &MockResult{
				Request: &r2.Request{
					Err: err,
				},
			}
		}
	}

	result := &MockResult{
		App: app,
		Request: &r2.Request{
			Request: req,
		},
	}

	// if the app isn't already started.
	if !app.IsStarted() {
		startupErrors := make(chan error)
		app.Config.BindAddr = DefaultMockBindAddr
		app.Config.ShutdownGracePeriod = time.Millisecond

		// set the on response delegate to stop the app.
		result.Request.Closer = func() error {
			return result.Close()
		}

		go func() {
			if err := app.Start(); err != nil {
				startupErrors <- err
			}
		}()

		select {
		case <-app.NotifyStarted():
		case err := <-startupErrors:
			result.Err = err
		}
	}

	if result.Request.URL == nil {
		result.Request.URL = &url.URL{
			Scheme: SchemeHTTP,
		}
	}
	if app.TLSConfig != nil {
		result.Request.URL.Scheme = SchemeHTTPS
	} else {
		result.Request.URL.Scheme = SchemeHTTP
	}

	if app.Listener == nil {
		result.Err = errors.New("the app listener is unset")
		return result
	}

	result.Request.URL.Host = app.Listener.Addr().String()

	return result
}

// MockMethod sends a mock request with a given method to an app.
// You should use request options to set the body of the request if it's a post or put etc.
func MockMethod(app *App, method, path string, options ...webutil.RequestOption) *MockResult {
	req := &http.Request{
		Method: method,
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockGet sends a mock get request to an app.
func MockGet(app *App, path string, options ...webutil.RequestOption) *MockResult {
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockPost sends a mock post request to an app.
func MockPost(app *App, path string, body io.ReadCloser, options ...webutil.RequestOption) *MockResult {
	req := &http.Request{
		Method: "POST",
		Body:   body,
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockResult is a result of a mocked request.
type MockResult struct {
	*r2.Request
	App *App
}

// Close stops the app.
func (mr *MockResult) Close() error {
	if mr.App.CanStop() {
		if err := mr.App.Stop(); err != nil {
			return err
		}
		<-mr.App.NotifyStopped()
	}
	return nil
}

// MockCtx returns a new mock ctx.
// It is intended to be used in testing.
func MockCtx(method, path string, options ...CtxOption) *Ctx {
	return NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest(method, path), options...)
}
