package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAsFlags(t *testing.T) {
	assert := assert.New(t)

	assert.Len(3, AsStrings(Error, Warning, Fatal))
	assert.Len(3, AsFlags("Error", "Warning", "Fatal"))

	flags := []string{"foo", "bar", "baz", "buzz"}
	assert.Equal(flags, AsStrings(AsFlags(flags...)...))
}
