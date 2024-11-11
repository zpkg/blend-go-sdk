/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

import (
	"context"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/configutil"
	"github.com/zpkg/blend-go-sdk/env"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("OAUTH_REDIRECT_URI", "https://app.com/oauth/google")
	env.Env().Set("OAUTH_HOSTED_DOMAIN", "foo.com")
	env.Env().Set("OAUTH_CLIENT_ID", "foo")
	env.Env().Set("OAUTH_CLIENT_SECRET", "bar")

	cfg := &Config{}
	ctx := env.WithVars(context.Background(), env.Env())
	assert.Nil(cfg.Resolve(ctx))
	assert.Equal("foo", cfg.ClientID)
	assert.Equal("bar", cfg.ClientSecret)
	assert.Equal("https://app.com/oauth/google", cfg.RedirectURI)
	assert.Equal("foo.com", cfg.HostedDomain)
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config{}.IsZero())
	assert.True(Config{ClientID: "foo"}.IsZero())
	assert.False(Config{ClientID: "foo", ClientSecret: "bar"}.IsZero())
}
