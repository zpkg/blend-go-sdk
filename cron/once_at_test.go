/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func Test_OnceAtUTC(t *testing.T) {
	assert := assert.New(t)

	fireAt := time.Date(2018, 10, 21, 12, 00, 00, 00, time.UTC)
	before := fireAt.Add(-time.Minute)
	after := fireAt.Add(time.Minute)

	s := OnceAtUTC(fireAt)
	result := s.Next(before)
	assert.Equal(result, fireAt)

	result = s.Next(after)
	assert.True(result.IsZero())
}

func Test_OnceAtUTC_String(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Now().UTC()

	its.Equal(
		fmt.Sprintf("%s %s", StringScheduleOnceAt, ts.Format(time.RFC3339)),
		OnceAtUTC(ts).String(),
	)
}
