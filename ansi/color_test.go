package ansi

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestColorApply(t *testing.T) {
	assert := assert.New(t)

	escapedBlack := ColorBlack.Normal()
	assert.Equal("\033[0;"+string(ColorBlack), escapedBlack)

	appliedBlack := ColorBlack.Apply("test")
	assert.Equal(ColorBlack.Normal()+"test"+ColorReset, appliedBlack)
}
