/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import (
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestToFloat64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1550059200000000000, ToFloat64(time.Date(2019, 02, 13, 12, 0, 0, 0, time.UTC)))
}
