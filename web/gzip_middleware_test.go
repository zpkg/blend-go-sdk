/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/r2"
	"github.com/zpkg/blend-go-sdk/webutil"
)

func TestGZipMiddlewarePlaintext(t *testing.T) {
	assert := assert.New(t)

	app := MustNew(
		OptUse(GZip),
	)
	app.GET("/", ok)

	req := MockGet(app, "/")
	assert.Nil(req.Err)
	resBody, _, err := req.Bytes()
	assert.Nil(err)
	assert.Equal("\"OK!\"\n", string(resBody))
}

func TestGZipMiddlewareCompressed(t *testing.T) {
	assert := assert.New(t)

	app := MustNew(
		OptUse(GZip),
	)
	app.GET("/", ok)

	req := MockGet(app, "/", r2.OptHeaderValue(webutil.HeaderAcceptEncoding, "gzip"))
	assert.Nil(req.Err)
	body, meta, err := req.Bytes()

	assert.Equal("gzip", meta.Header.Get(webutil.HeaderContentEncoding))
	assert.Equal(webutil.HeaderAcceptEncoding, meta.Header.Get(webutil.HeaderVary))
	assert.Nil(err)

	decompressor, err := gzip.NewReader(bytes.NewBuffer(body))
	assert.Nil(err)
	decompressed, err := io.ReadAll(decompressor)
	assert.Nil(err)

	assert.Equal("\"OK!\"\n", string(decompressed))
}
