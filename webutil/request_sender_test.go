package webutil

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewRequestSender(t *testing.T) {
	assert := assert.New(t)

	rs := NewRequestSender(MustParseURL("https://foo.com/bar"))
	assert.NotNil(rs.Transport())
	assert.NotNil(rs.Headers())
	assert.Equal(DefaultRequestTimeout, rs.Client().Timeout)

	assert.NotEqual("bar", rs.Headers().Get("foo"))
	rs.WithHeader("foo", "bar")
	assert.Equal("bar", rs.Headers().Get("foo"))

	assert.Equal(DefaultRequestMethod, rs.Method())
	rs.WithMethod("GET")
	assert.Equal("GET", rs.Method())

	assert.Nil(rs.Tracer())
	rs.WithTracer(mockRequestTracer{})
	assert.NotNil(rs.Tracer())
}

type mockRequestTracer struct {
	OnStart  func(*http.Request)
	OnFinish func(*http.Request, *http.Response, error)
}

func (mrt mockRequestTracer) Start(req *http.Request) RequestTraceFinisher {
	if mrt.OnStart != nil {
		mrt.OnStart(req)
	}
	return mockRequestTraceFinisher{Parent: &mrt}
}

type mockRequestTraceFinisher struct {
	Parent *mockRequestTracer
}

func (mrtf mockRequestTraceFinisher) Finish(req *http.Request, res *http.Response, err error) {
	if mrtf.Parent.OnFinish != nil {
		mrtf.Parent.OnFinish(req, res, err)
	}
}

func TestRequestSenderSend(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != DefaultRequestMethod {
			http.Error(w, "wrong method", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!\n")
	}))
	defer ts.Close()

	rs := NewRequestSender(MustParseURL(ts.URL))
	res, err := rs.Send()
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestSenderSendMethod(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "wrong method", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!\n")
	}))
	defer ts.Close()

	rs := NewRequestSender(MustParseURL(ts.URL)).WithMethod("GET")
	res, err := rs.Send()
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestSenderSendBytes(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if string(contents) != "foo bar baz" {
			http.Error(w, "wrong post body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!\n")
	}))
	defer ts.Close()

	rs := NewRequestSender(MustParseURL(ts.URL))
	res, err := rs.SendBytes(context.TODO(), []byte("foo bar baz"))
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestSenderSendJSON(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if string(contents) != "{\"foo\":\"bar\"}" {
			http.Error(w, "wrong post body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!\n")
	}))
	defer ts.Close()

	rs := NewRequestSender(MustParseURL(ts.URL))
	res, err := rs.SendJSON(context.TODO(), map[string]interface{}{"foo": "bar"})
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func TestRequestSenderSendTraced(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != DefaultRequestMethod {
			http.Error(w, "wrong method", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!\n")
	}))
	defer ts.Close()

	var didCallStart, didCallFinish bool
	tracer := mockRequestTracer{
		OnStart: func(req *http.Request) {
			assert.NotNil(req)
			didCallStart = true
		},
		OnFinish: func(req *http.Request, res *http.Response, err error) {
			didCallFinish = true
		},
	}

	rs := NewRequestSender(MustParseURL(ts.URL)).WithTracer(tracer)
	res, err := rs.Send()
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal(http.StatusOK, res.StatusCode)

	assert.True(didCallStart)
	assert.True(didCallFinish)
}
