package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/blend/go-sdk/exception"
)

// Decompress gzip decompresses the bytes.
func Decompress(contents []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(contents))
	if err != nil {
		return nil, exception.New(err)
	}
	defer r.Close()
	decompressed, err := ioutil.ReadAll(r)
	return decompressed, exception.New(err)
}
