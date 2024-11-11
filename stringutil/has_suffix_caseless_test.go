/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestHasSuffixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.True(HasSuffixCaseless("hello world!", "world!"))
	assert.True(HasSuffixCaseless("hello world", "d"))
	assert.True(HasSuffixCaseless("hello world", "hello world"))

	assert.True(HasSuffixCaseless("hello WORLD", "world"))
	assert.True(HasSuffixCaseless("hello world", "WORLD"))

	assert.False(HasSuffixCaseless("hello world", "hello hello world"))
	assert.False(HasSuffixCaseless("hello world", "foobar"))
	assert.False(HasSuffixCaseless("hello world", "hello")) //this would pass prefix
}
