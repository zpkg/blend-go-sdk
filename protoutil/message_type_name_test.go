/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package protoutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/protoutil/testdata"
)

func Test_MessageTypeName(t *testing.T) {
	its := assert.New(t)

	its.Equal("testdata.Message", MessageTypeName(new(testdata.Message)))
}
