package proxy

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestUpstreamWithoutHopHeaders(t *testing.T) {
	assert := assert.New(t)

	u := &Upstream{HopHeaders: []string{"foo", "bar", "baz", "buzz"}}
	u = u.WithoutHopHeaders("bar", "baz")
	assert.Len(2, u.HopHeaders)
	assert.Equal("foo", u.HopHeaders[0])
	assert.Equal("buzz", u.HopHeaders[1])
}
