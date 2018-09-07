package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStateValue(t *testing.T) {
	assert := assert.New(t)

	state := State{
		"foo":  "bar",
		"buzz": "fuzz",
	}

	assert.Equal("bar", state.Value("foo"))
	assert.Equal("fuzz", state.Value("buzz"))
}
