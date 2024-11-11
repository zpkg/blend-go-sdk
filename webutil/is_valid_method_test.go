/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestIsValidMethod(t *testing.T) {
	assert := assert.New(t)

	methods := []string{
		MethodGet,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodOptions,
	}

	for _, method := range methods {
		assert.True(IsValidMethod(method))
	}

	assert.False(IsValidMethod("\n"))
}
