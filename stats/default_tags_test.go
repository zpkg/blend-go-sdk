/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/env"
)

func TestAddDefaultTagsFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("SERVICE_NAME", "someservice")
	env.Env().Set("SERVICE_ENV", "sandbox")
	env.Env().Set("HOSTNAME", "somecontainer")

	// Handles nil collector
	AddDefaultTagsFromEnv(nil)

	collector := NewMockCollector(32)
	AddDefaultTagsFromEnv(collector)

	tags := collector.DefaultTags()
	assert.Len(tags, 3)
	assert.Equal("service:someservice", tags[0])
	assert.Equal("env:sandbox", tags[1])
	assert.Equal("container:somecontainer", tags[2])
}

func TestAddDefaultTags(t *testing.T) {
	assert := assert.New(t)

	// Handles nil collector
	AddDefaultTagsFromEnv(nil)

	collector := NewMockCollector(32)
	AddDefaultTags(collector, "someservice", "sandbox", "somecontainer")

	tags := collector.DefaultTags()
	assert.Len(tags, 3)
	assert.Equal("service:someservice", tags[0])
	assert.Equal("env:sandbox", tags[1])
	assert.Equal("container:somecontainer", tags[2])
}
