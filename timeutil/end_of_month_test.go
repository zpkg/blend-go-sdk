/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestEndOfMonth(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    time.Time
		Expected time.Time
	}{
		{Input: time.Date(2019, 9, 9, 17, 59, 44, 0, time.UTC), Expected: time.Date(2019, 9, 30, 23, 59, 59, 0, time.UTC)},
		{Input: time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC), Expected: time.Date(2019, 9, 30, 23, 59, 59, 0, time.UTC)},
		{Input: time.Date(2019, 9, 30, 23, 59, 59, 0, time.UTC), Expected: time.Date(2019, 9, 30, 23, 59, 59, 0, time.UTC)},
	}

	for _, tc := range testCases {
		assert.InTimeDelta(
			tc.Expected,
			EndOfMonth(tc.Input),
			time.Second,
			fmt.Sprintf("input: %v expected: %v", tc.Input, tc.Expected),
		)
	}
}
