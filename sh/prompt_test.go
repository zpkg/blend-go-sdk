package sh

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPromptFrom(t *testing.T) {
	assert := assert.New(t)

	input := bytes.NewBufferString("test\n")
	output := bytes.NewBuffer(nil)

	value, err := PromptFrom(output, input, "value: ")
	assert.Nil(err)
	assert.Equal("test", value)
}

func TestPromptFromEmpty(t *testing.T) {
	assert := assert.New(t)

	input := bytes.NewBufferString("\n")
	output := bytes.NewBuffer(nil)

	value, err := PromptFrom(output, input, "value: ")
	assert.Nil(err)
	assert.Equal("", value)
}
