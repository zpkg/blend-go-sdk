/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// AddListeners adds error listeners.
func AddListeners(log logger.Listenable, meta configmeta.Meta, cfg Config) error {
	if log == nil || cfg.IsZero() {
		return nil
	}
	if typed, ok := log.(logger.InfofReceiver); ok {
		typed.Infof("using sentry host: %s", webutil.MustParseURL(cfg.DSN).Hostname())
	}
	if cfg.Environment == "" {
		cfg.Environment = meta.ServiceEnv
	}
	if cfg.ServerName == "" {
		cfg.ServerName = meta.ServiceName
	}
	if cfg.Release == "" {
		cfg.Release = meta.Version
	}
	client, err := New(cfg)
	if err != nil {
		return err
	}
	listener := logger.NewErrorEventListener(client.Notify)
	log.Listen(logger.Error, ListenerName, listener)
	log.Listen(logger.Fatal, ListenerName, listener)
	return nil
}
