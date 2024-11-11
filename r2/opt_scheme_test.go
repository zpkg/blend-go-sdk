/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestOptScheme(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptScheme("spdy"))
	assert.Equal("spdy", r.Request.URL.Scheme)

	var unset Request
	assert.NotNil(OptScheme("spdy")(&unset))
}
