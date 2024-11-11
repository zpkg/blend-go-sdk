/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"time"

	"github.com/zpkg/blend-go-sdk/db"
	"github.com/zpkg/blend-go-sdk/ex"
)

func verifyStatementTimeout(ctx context.Context, pool *db.Connection) (time.Duration, error) {
	type statementTimeoutRow struct {
		StatementTimeout string `db:"statement_timeout"`
	}

	q := pool.QueryContext(ctx, "SHOW statement_timeout;")
	r := statementTimeoutRow{}
	found, err := q.Out(&r)
	if !found {
		return 0, ex.New("`SHOW statement_timeout;` query returned no results")
	}
	if err != nil {
		return 0, err
	}

	d, err := time.ParseDuration(r.StatementTimeout)
	if err != nil {
		return 0, err
	}

	return d, nil
}

func ensureStatementTimeout(ctx context.Context, pool *db.Connection, cfg *config) (time.Duration, error) {
	timeout, err := verifyStatementTimeout(ctx, pool)
	if err != nil {
		return 0, err
	}

	if timeout != cfg.StatementTimeout {
		err = ex.New(
			"Expected the default statement timeout to be set",
			ex.OptMessagef("Timeout: %s, Expected: %s", timeout, cfg.StatementTimeout),
		)
		return 0, err
	}

	return timeout, nil
}
