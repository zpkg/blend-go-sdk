/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"encoding/json"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestClassMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	err := Class("this is only a test")
	contents, marshalErr := json.Marshal(err)
	assert.Nil(marshalErr)
	assert.Equal(`"this is only a test"`, string(contents))
}
