package ansi

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestColorApply(t *testing.T) {
	assert := assert.New(t)

	escapedBlack := ColorBlack.Escaped()
	assert.Equal("\033["+string(ColorBlack), escapedBlack)

	appliedBlack := ColorBlack.Apply("test")
	assert.Equal(ColorBlack.Escaped()+"test"+ColorReset.Escaped(), appliedBlack)
}
