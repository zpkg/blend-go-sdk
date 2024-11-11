/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func Test_SchemeIsSecure(t *testing.T) {
	its := assert.New(t)

	its.True(SchemeIsSecure(SchemeHTTPS))
	its.True(SchemeIsSecure(SchemeSPDY))

	its.False(SchemeIsSecure(SchemeHTTP))
	its.False(SchemeIsSecure("garbage"))
	its.False(SchemeIsSecure(""))
}
