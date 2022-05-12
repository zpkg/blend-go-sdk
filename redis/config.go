/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis

import (
	"context"
	"time"

	"github.com/blend/go-sdk/configutil"
)

// Config is the config type for the redis client.
type Config struct {
	Network        string        `yaml:"network"`
	Addr           string        `yaml:"addr"`
	AuthUser       string        `yaml:"authUser"`
	AuthPassword   string        `yaml:"authPassword"`
	DB             string        `yaml:"db"`
	ConnectTimeout time.Duration `yaml:"connectTimeout"`
	Timeout        time.Duration `yaml:"timeout"`
}

// Resolve resolves the config.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.Network, configutil.Env("REDIS_NETWORK"), configutil.String(c.Network), configutil.String(DefaultNetwork)),
		configutil.SetString(&c.Addr, configutil.Env("REDIS_ADDR"), configutil.String(c.Addr), configutil.String(DefaultAddr)),
		configutil.SetString(&c.DB, configutil.Env("REDIS_DB"), configutil.String(c.DB)),
		configutil.SetString(&c.AuthUser, configutil.Env("REDIS_AUTH_USER"), configutil.String(c.AuthUser)),
		configutil.SetString(&c.AuthPassword, configutil.Env("REDIS_AUTH_PASS"), configutil.String(c.AuthPassword)),
		configutil.SetDuration(&c.ConnectTimeout, configutil.Env("REDIS_CONNECT_TIMEOUT"), configutil.Duration(c.ConnectTimeout), configutil.Duration(DefaultConnectTimeout)),
		configutil.SetDuration(&c.Timeout, configutil.Env("REDIS_TIMEOUT"), configutil.Duration(c.Timeout), configutil.Duration(DefaultTimeout)),
	)
}
