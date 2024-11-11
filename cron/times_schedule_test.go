/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestTimesSchedule(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Date(2020, 01, 23, 12, 11, 10, 9, time.UTC)
	schedule := Times(5, Every(time.Millisecond))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.True(schedule.Next(ts).IsZero())
}
