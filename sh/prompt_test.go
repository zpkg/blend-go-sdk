/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sh

import (
	"bytes"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestPromptFrom(t *testing.T) {
	assert := assert.New(t)

	input := bytes.NewBufferString("test\n")
	output := bytes.NewBuffer(nil)
	assert.Equal("test", PromptFrom(output, input, "value: "))
}
