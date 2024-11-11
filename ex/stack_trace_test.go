/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"fmt"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestGetStackTrace(t *testing.T) {
	assert := assert.New(t)

	assert.NotEmpty(GetStackTrace())
}

func TestStackStrings(t *testing.T) {
	assert := assert.New(t)

	stack := []string{
		"foo",
		"bar",
		"baz",
	}

	stackStrings := StackStrings(stack)

	assert.Equal("\nfoo\nbar\nbaz", fmt.Sprintf("%+v", stackStrings))
	assert.Equal("[]string{\"foo\", \"bar\", \"baz\"}", fmt.Sprintf("%#v", stackStrings))
	assert.Equal("\nfoo\nbar\nbaz", fmt.Sprintf("%v", stackStrings))
	assert.Equal([]string{"foo", "bar", "baz"}, stackStrings)
}

func TestExceptionWithStackStrings(t *testing.T) {
	assert := assert.New(t)

	stack := []string{
		"foo",
		"bar",
		"baz",
	}

	ex := As(New("foo", OptStackTrace(StackStrings(stack))))

	values := ex.Decompose()
	assert.NotEmpty(values["StackTrace"])
	assert.NotNil(ex.StackTrace)
}
