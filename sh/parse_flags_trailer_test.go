package sh

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestParseFlagsTrailer(t *testing.T) {
	assert := assert.New(t)

	parsed, err := ParseFlagsTrailer("foo", "bar")
	assert.True(exception.Is(err, ErrFlagsNoTrailer))
	assert.Empty(parsed)

	parsed, err = ParseFlagsTrailer("foo", "bar", "--")
	assert.True(exception.Is(err, ErrFlagsNoTrailer))
	assert.Empty(parsed)

	parsed, err = ParseFlagsTrailer("foo", "bar", "--", "echo", "'things'")
	assert.Nil(err)
	assert.Equal("echo 'things'", parsed)
}
