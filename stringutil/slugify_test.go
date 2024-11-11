/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestSlugify(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input, Expected string
	}{
		{"foo", "foo"},
		{"Foo", "foo"},
		{"f00", "f00"},
		{"foo-bar", "foo-bar"},
		{"foo & bar", "foo-bar"},
		{"foo--bar", "foo-bar"},
		{"foo-.bar", "foo-bar"},
		{"foo bar", "foo-bar"},
		{"foo  bar", "foo-bar"},
		{"foo\tbar", "foo-bar"},
		{"foo\nbar", "foo-bar"},
		{"foo\t\nbar", "foo-bar"},
		{"foo\t\nbar\t\n", "foo-bar-"},
		{"Mt. Tam", "mt-tam"},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, Slugify(tc.Input))
	}
}
