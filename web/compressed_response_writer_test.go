package web

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestBufferedCompressedWriter(t *testing.T) {
	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)
	mockedWriter := webutil.NewMockResponse(buf)
	bufferedWriter := NewCompressedResponseWriter(mockedWriter)

	written, err := bufferedWriter.Write([]byte("ok"))
	assert.Nil(err)
	assert.NotZero(written)
}
