/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package httpstats

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/stats"
	"github.com/zpkg/blend-go-sdk/webutil"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(webutil.FlagHTTPRequest, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector(32))
	assert.True(log.HasListener(webutil.FlagHTTPRequest, stats.ListenerNameStats))
}
