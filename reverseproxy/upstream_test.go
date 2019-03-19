package reverseproxy

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestUpstreamWithoutHopHeaders(t *testing.T) {
	assert := assert.New(t)

	u := NewUpstream(MustParseURL("http://localhost:5000"))
	assert.NotNil(u.ReverseProxy)
}
