/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestMock(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.GET("/", func(_ *Ctx) Result { return NoContent })

	res, err := Mock(app, &http.Request{Method: "GET", URL: &url.URL{Scheme: webutil.SchemeHTTP, Path: "/"}}).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)

	assert.True(app.IsStopped())
}

func TestMockGet(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()
	app.GET("/", func(_ *Ctx) Result { return NoContent })

	res, err := MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNoContent, res.StatusCode)

	assert.True(app.IsStopped())
}
