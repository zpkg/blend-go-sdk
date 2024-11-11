/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func Test_FileParseSize(t *testing.T) {
	assert := assert.New(t)

	parsed, err := ParseFileSize("2gb")
	assert.Nil(err)
	assert.Equal(2*Gigabyte, parsed)

	parsed, err = ParseFileSize("3mb")
	assert.Nil(err)
	assert.Equal(3*Megabyte, parsed)

	parsed, err = ParseFileSize("123kb")
	assert.Nil(err)
	assert.Equal(123*Kilobyte, parsed)

	parsed, err = ParseFileSize("12345")
	assert.Nil(err)
	assert.Equal(12345, parsed)

	parsed, err = ParseFileSize("")
	assert.Nil(err)
	assert.Equal(0, parsed)

	parsed, err = ParseFileSize("bogus")
	assert.NotNil(err)
	assert.Equal(0, parsed)
}
