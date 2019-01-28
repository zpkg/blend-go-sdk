package sh

import (
	"bytes"
	"io"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTeeWriter(t *testing.T) {
	assert := assert.New(t)

	buf0 := bytes.NewBuffer(nil)
	buf1 := bytes.NewBuffer(nil)

	input := bytes.NewBuffer([]byte(`this is only a test`))

	count, err := io.Copy(Tee(buf0, buf1), input)
	assert.Nil(err)
	assert.Equal(19, count)

	assert.Equal([]byte(`this is only a test`), buf0.Bytes())
	assert.Equal([]byte(`this is only a test`), buf1.Bytes())
}
