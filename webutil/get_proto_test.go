package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

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
