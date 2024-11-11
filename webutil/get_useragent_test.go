/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestGetUseragent(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("go-sdk test", GetUserAgent(NewMockRequest("GET", "/")))
}
