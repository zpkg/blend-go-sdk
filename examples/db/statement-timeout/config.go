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
	envVarStatementTimeout = "DB_STATEMENT_TIMEOUT"
)

type config struct {
	StatementTimeout time.Duration
	PGSleep          time.Duration
	ContextTimeout   time.Duration
}

func (c *config) SetEnvironment() error {
	existing := env.Env().String(envVarStatementTimeout)
	if existing != "" {
		err := ex.New(
			"Statement timeout will be set by the statement-timeout script",
			ex.OptMessagef("Value set from the environment, %s=%q", envVarStatementTimeout, existing),
		)
		return err
	}

	env.Env().Set(envVarStatementTimeout, c.StatementTimeout.String())
	return nil
}

func (c *config) Print() {
	if c == nil {
		return
	}

	log.Printf("Configured statement timeout: %s\n", c.StatementTimeout)
	log.Printf("Configured pg_sleep:          %s\n", c.PGSleep)
	log.Printf("Configured context timeout:   %s\n", c.ContextTimeout)
}

func getConfig() *config {
	viaGoContext := env.Env().String("VIA_GO_CONTEXT")
	if strings.EqualFold(viaGoContext, "true") {
		return &config{
			StatementTimeout: 10 * time.Second,
			PGSleep:          200 * time.Millisecond,
			ContextTimeout:   100 * time.Millisecond,
		}
	}

	return &config{
		StatementTimeout: 10 * time.Millisecond,
		PGSleep:          200 * time.Millisecond,
		ContextTimeout:   400 * time.Millisecond,
	}
}
