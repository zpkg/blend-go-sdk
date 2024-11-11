/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"log"
	"strings"
	"time"

	"github.com/zpkg/blend-go-sdk/env"
	"github.com/zpkg/blend-go-sdk/ex"
)

const (
	envVarLockTimeout = "DB_LOCK_TIMEOUT"
)

type config struct {
	LockTimeout    time.Duration
	ContextTimeout time.Duration
	TxSleep        time.Duration
}

func (c *config) SetEnvironment() error {
	existing := env.Env().String(envVarLockTimeout)
	if existing != "" {
		err := ex.New(
			"Lock timeout will be set by the prevent-deadlock script",
			ex.OptMessagef("Value set from the environment, %s=%q", envVarLockTimeout, existing),
		)
		return err
	}

	env.Env().Set(envVarLockTimeout, c.LockTimeout.String())
	return nil
}

func (c *config) Print() {
	if c == nil {
		return
	}

	log.Printf("Configured lock timeout:      %s\n", c.LockTimeout)
	log.Printf("Configured context timeout:   %s\n", c.ContextTimeout)
	log.Printf("Configured transaction sleep: %s\n", c.TxSleep)
}

func getConfig() *config {
	forceDeadlock := env.Env().String("FORCE_DEADLOCK")
	if strings.EqualFold(forceDeadlock, "true") {
		return &config{
			LockTimeout:    10 * time.Second,
			ContextTimeout: 10 * time.Second,
			TxSleep:        200 * time.Millisecond,
		}
	}

	betweenQueries := env.Env().String("BETWEEN_QUERIES")
	if strings.EqualFold(betweenQueries, "true") {
		return &config{
			LockTimeout:    10 * time.Second,
			ContextTimeout: 100 * time.Millisecond,
			TxSleep:        200 * time.Millisecond,
		}
	}

	disable := env.Env().String("DISABLE_LOCK_TIMEOUT")
	if strings.EqualFold(disable, "true") {
		return &config{
			LockTimeout:    10 * time.Second,
			ContextTimeout: 600 * time.Millisecond,
			TxSleep:        200 * time.Millisecond,
		}
	}

	return &config{
		LockTimeout:    10 * time.Millisecond,
		ContextTimeout: 600 * time.Millisecond,
		TxSleep:        200 * time.Millisecond,
	}
}
