package fileutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFileReadByLines(t *testing.T) {
	assert := assert.New(t)

	called := false
	ReadLines("README.md", func(line string) error {
		called = true
		return nil
	})

	assert.True(called, "We should have called the handler for `README.md`")
}

func TestFileReadByChunks(t *testing.T) {
	assert := assert.New(t)

	called := false
	ReadChunks("README.md", 32, func(chunk []byte) error {
		called = true
		return nil
	})

	assert.True(called, "We should have called the handler for `README.md`")
}

func TestFileParseSize(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(2*Gigabyte, ParseFileSize("2gb", 1))
	assert.Equal(3*Megabyte, ParseFileSize("3mb", 1))
	assert.Equal(123*Kilobyte, ParseFileSize("123kb", 1))
	assert.Equal(12345, ParseFileSize("12345", 1))
	assert.Equal(1, ParseFileSize("", 1))
}
