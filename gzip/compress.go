package gzip

import (
	"bytes"
	"compress/gzip"

	"github.com/blend/go-sdk/exception"
)

// Compress gzip compresses the bytes.
func Compress(contents []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(contents)
	err := w.Flush()
	if err != nil {
		return nil, exception.New(err)
	}
	err = w.Close()
	if err != nil {
		return nil, exception.New(err)
	}

	return b.Bytes(), nil
}
