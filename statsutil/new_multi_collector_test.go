/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package statsutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
)

func Test_NewMultiCollector(t *testing.T) {
	its := assert.New(t)

	log := logger.None()

	collector, err := NewMultiCollector(log,
		OptMetaConfig(configmeta.Meta{
			ServiceName: "test-service",
			ServiceEnv:  "test-service-env",
			Version:     "test-service-version",
			Hostname:    "test-service-hostname",
		}),
		OptDatadogConfig(datadog.Config{}),
		OptPrinter(true),
	)
	its.Nil(err)

	typed, ok := collector.(stats.MultiCollector)

	its.True(ok)
	its.Len(typed, 2)
	its.True(typed.HasTagKey(stats.TagService))
	its.True(typed.HasTagKey(stats.TagEnv))
	its.True(typed.HasTagKey(stats.TagHostname))
	its.True(typed.HasTagKey(stats.TagVersion))

	defaultTags := typed.DefaultTags()
	its.Any(defaultTags, func(v interface{}) bool {
		return v.(string) == stats.Tag(stats.TagService, "test-service")
	})
	its.Any(defaultTags, func(v interface{}) bool {
		return v.(string) == stats.Tag(stats.TagEnv, "test-service-env")
	})
	its.Any(defaultTags, func(v interface{}) bool {
		return v.(string) == stats.Tag(stats.TagHostname, "test-service-hostname")
	})
	its.Any(defaultTags, func(v interface{}) bool {
		return v.(string) == stats.Tag(stats.TagVersion, "test-service-version")
	})
}
