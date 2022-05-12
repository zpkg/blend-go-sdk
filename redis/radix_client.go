/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis

import (
	"context"
	"time"

	"github.com/mediocregopher/radix/v4"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// New returns a new client.
func New(ctx context.Context, opts ...Option) (*RadixClient, error) {
	var rc RadixClient
	var err error
	for _, opt := range opts {
		if err = opt(&rc); err != nil {
			return nil, ex.New(err)
		}
	}
	if rc.Config.ConnectTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, rc.Config.ConnectTimeout)
		defer cancel()
	}
	rc.Client, err = (radix.PoolConfig{
		Dialer: radix.Dialer{
			SelectDB: rc.Config.DB,
			AuthUser: rc.Config.AuthUser,
			AuthPass: rc.Config.AuthPassword,
		},
	}).New(ctx, rc.Config.Network, rc.Config.Addr)
	if err != nil {
		return nil, ex.New(err)
	}
	return &rc, nil
}

// Assert `RadixClient` implements `Client`.
var (
	_ Client = (*RadixClient)(nil)
)

// RadixClient is a wrapping client for the underling radix redis driver.
type RadixClient struct {
	Config Config
	Log    logger.Triggerable
	Tracer Tracer
	Client radix.Client
}

// Do runs a given command.
func (rc *RadixClient) Do(ctx context.Context, out interface{}, op string, args ...string) (err error) {
	if rc.Log != nil {
		started := time.Now()
		defer func() {
			rc.Log.TriggerContext(ctx, NewEvent(op, args, time.Since(started),
				OptEventNetwork(rc.Config.Network),
				OptEventAddr(rc.Config.Addr),
				OptEventAuthUser(rc.Config.AuthUser),
				OptEventDB(rc.Config.DB),
				OptEventErr(err),
			))
		}()
	}
	if rc.Tracer != nil {
		finisher := rc.Tracer.Do(ctx, rc.Config, op, args)
		defer finisher.Finish(ctx, err)
	}
	if rc.Config.Timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, rc.Config.Timeout)
		defer cancel()
	}
	if radixErr := rc.Client.Do(ctx, radix.Cmd(out, op, args...)); radixErr != nil {
		err = ex.New(radixErr)
		return
	}
	return
}

// Close closes the underlying connection.
func (rc *RadixClient) Close() error {
	return rc.Client.Close()
}
