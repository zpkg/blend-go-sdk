/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ref

import (
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestRef(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(String("foo"))
	assert.NotEmpty(Strings("foo", "bar"))

	assert.NotNil(Bool(true))

	assert.NotNil(Byte('b'))
	assert.NotNil(Rune('b'))

	assert.NotNil(Uint8(0))
	assert.NotNil(Uint16(0))
	assert.NotNil(Uint32(0))
	assert.NotNil(Uint64(0))
	assert.NotNil(Int8(0))
	assert.NotNil(Int16(0))
	assert.NotNil(Int32(0))
	assert.NotNil(Int64(0))
	assert.NotNil(Int(0))
	assert.NotNil(Float32(0))
	assert.NotNil(Float64(0))
	assert.NotNil(Time(time.Time{}))
	assert.NotNil(Duration(0))
}
