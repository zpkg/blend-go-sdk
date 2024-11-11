/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2test

import (
	"net/http"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/r2"
)

func TestOptMockResponseString(t *testing.T) {
	it := assert.New(t)

	var didCallOriginalCloser bool
	output, meta, err := r2.New(r2.TestURL,
		r2.OptPost(),
		r2.OptCloser(func() error {
			didCallOriginalCloser = true
			return nil
		}),
		OptMockResponseString("this is just a test!"),
	).Bytes()
	it.Nil(err)
	it.Equal(http.StatusOK, meta.StatusCode)
	it.Equal("this is just a test!", string(output))
	it.True(didCallOriginalCloser)
}

func TestOptMockResponseStringStatus(t *testing.T) {
	it := assert.New(t)

	var didCallOriginalCloser bool
	output, meta, err := r2.New(r2.TestURL,
		r2.OptPost(),
		r2.OptCloser(func() error {
			didCallOriginalCloser = true
			return nil
		}),
		OptMockResponseStringStatus(http.StatusForbidden, "this is just a test!"),
	).Bytes()
	it.Nil(err)
	it.Equal(http.StatusForbidden, meta.StatusCode)
	it.Equal("this is just a test!", string(output))
	it.True(didCallOriginalCloser)
}
