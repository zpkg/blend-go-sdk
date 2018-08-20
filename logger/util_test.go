package logger

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	assert := assert.New(t)

	hdr := http.Header{}
	hdr.Set("X-Forwarded-For", "123")
	r := http.Request{
		Header: hdr,
	}
	assert.Equal("123", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-FORWARDED-FOR", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-FORWARDED-FOR", "1,2,3")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-Real-Ip", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-REAL-IP", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	hdr = http.Header{}
	hdr.Set("X-REAL-IP", "1,2,3")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetRemoteAddr(&r))

	r = http.Request{
		RemoteAddr: "1:1",
	}
	assert.Equal("1", GetRemoteAddr(&r))

	r = http.Request{
		RemoteAddr: "1",
	}
	assert.Equal("", GetRemoteAddr(&r))
}

func TestGetProto(t *testing.T) {
	assert := assert.New(t)

	headers := http.Header{}
	headers.Set("X-Forwarded-Proto", "https")
	r := http.Request{
		Proto:  "http",
		Header: headers,
	}
	assert.Equal("https", GetProto(&r))

	headers = http.Header{}
	headers.Set("X-Forwarded-Proto", "spdy,https")
	r = http.Request{
		Proto:  "http",
		Header: headers,
	}
	assert.Equal("spdy", GetProto(&r))

	headers = http.Header{}
	r = http.Request{
		Proto:  "http",
		Header: headers,
	}
	assert.Equal("http", GetProto(&r))
}

func TestGetHost(t *testing.T) {
	assert := assert.New(t)

	r := http.Request{
		Host: "local.test.com",
	}
	assert.Equal("local.test.com", GetHost(&r))

	r = http.Request{
		Host: "local.test.com:8080",
	}
	assert.Equal("local.test.com", GetHost(&r))

	r = http.Request{
		URL:  &url.URL{},
		Host: "local.test.com:8080",
	}
	assert.Equal("local.test.com", GetHost(&r))

	r = http.Request{
		URL:  &url.URL{Host: "local.foo.com"},
		Host: "local.test.com:8080",
	}
	assert.Equal("local.foo.com", GetHost(&r))

	headers := http.Header{}
	headers.Set("X-Forwarded-Proto", "spdy,https")
	headers.Set("X-Forwarded-Host", "local.bar.com")
	r = http.Request{
		URL:    &url.URL{Host: "local.foo.com"},
		Host:   "local.test.com:8080",
		Header: headers,
	}
	assert.Equal("local.bar.com", GetHost(&r))
}
