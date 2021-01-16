/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewEvent(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(FlagComplete, "test_task")
	assert.Equal(FlagComplete, e.GetFlag())
	assert.Equal("test_task", e.JobName)
}
