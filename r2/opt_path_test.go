/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptPath(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptPath("/not-foo"))
	assert.Equal("/not-foo", r.Request.URL.Path)

	var unset Request
	assert.NotNil(OptPath("/not-foo")(&unset))
}

func TestOptPathf(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptPathf("/not-foo/%s", "bar"))
	assert.Equal("/not-foo/bar", r.Request.URL.Path)

	var unset Request
	assert.NotNil(OptPathf("/not-foo/%s", "bar")(&unset))
}
