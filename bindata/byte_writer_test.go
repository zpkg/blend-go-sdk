/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT
license that can be found in the LICENSE file.

*/

package bindata

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestByteWriter(t *testing.T) {
	assert := assert.New(t)

	contents := []byte(`this is a test`)
	buffer := new(bytes.Buffer)

	bw := NewByteWriter(buffer)
	n, err := bw.Write(contents)
	assert.Nil(err)
	assert.Equal(14, n)
	assert.NotEmpty(buffer.Bytes())
}
