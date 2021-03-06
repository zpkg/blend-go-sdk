/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/logger"
)

func TestAddListeners(t *testing.T) {
	its := assert.New(t)

	its.Nil(AddListeners(nil, configmeta.Meta{}, Config{}))

	log := logger.None()
	its.Nil(AddListeners(log, configmeta.Meta{}, Config{}))
	its.False(log.HasListeners(logger.Error))
	its.False(log.HasListeners(logger.Fatal))

	its.Nil(AddListeners(log, configmeta.Meta{}, Config{DSN: "http://foo@example.org/1"}))
	its.True(log.HasListeners(logger.Error))
	its.True(log.HasListeners(logger.Fatal))

	its.True(log.HasListener(logger.Error, ListenerName))
	its.True(log.HasListener(logger.Fatal, ListenerName))
}
