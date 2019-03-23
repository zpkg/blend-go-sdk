package logger

import (
	"bytes"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestHTTPRequestEventListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)

	textBuffer := bytes.NewBuffer(nil)
	jsonBuffer := bytes.NewBuffer(nil)
	all := New().WithFlags(AllFlags()).WithRecoverPanics(false).
		WithWriter(NewTextWriter(textBuffer)).
		WithWriter(NewJSONWriter(jsonBuffer))
	defer all.Close()

	all.Listen(HTTPRequest, "default", NewHTTPRequestEventListener(func(hre *HTTPRequestEvent) {
		defer wg.Done()
		assert.Equal(HTTPRequest, hre.Flag())
		assert.NotNil(hre.Request())
		assert.Equal("test.com", hre.Request().Host)
	}))

	go func() {
		defer wg.Done()
		all.Trigger(NewHTTPRequestEvent(&http.Request{Host: "test.com", URL: &url.URL{}}))
	}()
	go func() {
		defer wg.Done()
		all.Trigger(NewHTTPRequestEvent(&http.Request{Host: "test.com", URL: &url.URL{}}))
	}()
	wg.Wait()
	all.Drain()

	assert.NotEmpty(textBuffer.String())
	assert.NotEmpty(jsonBuffer.String())
}

func TestHTTPRequestEventProperties(t *testing.T) {
	assert := assert.New(t)

	e := NewHTTPRequestEvent(nil)

	assert.False(e.Timestamp().IsZero())
	assert.True(e.WithTimestamp(time.Time{}).Timestamp().IsZero())

	assert.Empty(e.Labels())
	assert.Equal("bar", e.WithLabel("foo", "bar").Labels()["foo"])

	assert.Empty(e.Annotations())
	assert.Equal("zar", e.WithAnnotation("moo", "zar").Annotations()["moo"])

	assert.Equal(HTTPRequest, e.Flag())
	assert.Equal(Error, e.WithFlag(Error).Flag())

	assert.Empty(e.Headings())
	assert.Equal([]string{"Heading"}, e.WithHeadings("Heading").Headings())

	assert.Nil(e.Request())
	assert.NotNil(e.WithRequest(&http.Request{}).Request())

	assert.Nil(e.State())
	assert.Equal("foo", e.WithState(map[interface{}]interface{}{"bar": "foo"}).State()["bar"])

	assert.Empty(e.Route())
	assert.Equal("Route", e.WithRoute("Route").Route())
}
