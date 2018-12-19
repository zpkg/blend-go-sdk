package gzip

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stringutil"
)

func TestCompress(t *testing.T) {
	assert := assert.New(t)

	contents := []byte(stringutil.Random(stringutil.Letters, 1024))
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
