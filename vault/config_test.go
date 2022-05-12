/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

func TestConfigIsZero(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config{}.IsZero())
	assert.False(Config{Token: "garbage"}.IsZero())
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	assert.Equal(DefaultAddr, cfg.AddrOrDefault())
	assert.Empty(cfg.Token)
	assert.Equal(DefaultTimeout, cfg.TimeoutOrDefault())
	assert.Empty(cfg.RootCAs)
}
