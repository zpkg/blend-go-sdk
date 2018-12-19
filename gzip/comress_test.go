package gzip

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCompress(t *testing.T) {
	assert := assert.New(t)

	contents := []byte(String.Random(1024))
	compressed, err := Compress(contents)
	assert.Nil(err)
	assert.NotEmpty(compressed)
	assert.NotEqual(string(compressed), string(contents))

	assert.True(len(contents) > len(compressed))

	decompressed, err := Decompress(compressed)
	assert.Nil(err)
	assert.NotEmpty(decompressed)
	assert.Equal(string(contents), string(decompressed))
}
