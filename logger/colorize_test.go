/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/ansi"
	"github.com/zpkg/blend-go-sdk/assert"
)

func TestFlagTextColor(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Flag     string
		Expected ansi.Color
	}{
		{Info, ansi.ColorLightWhite},
		{Debug, ansi.ColorLightYellow},
		{Warning, ansi.ColorLightYellow},
		{Error, ansi.ColorRed},
		{Fatal, ansi.ColorRed},
		{"foo", DefaultFlagTextColor},
		{"", DefaultFlagTextColor},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, FlagTextColor(tc.Flag))
	}
}
