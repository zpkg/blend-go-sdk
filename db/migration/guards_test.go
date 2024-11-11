/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/db"
)

func TestGuard(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer func() { _ = tx.Rollback() }()

	tableName := randomName()
	err = createTestTable(tableName, tx)
	assert.Nil(err)

	err = insertTestValue(tableName, 4, "test", tx)
	assert.Nil(err)

	var didRun bool
	action := Actions(ActionFunc(func(ctx context.Context, c *db.Connection, itx *sql.Tx) error {
		didRun = true
		return nil
	}))

	err = Guard("test", func(ctx context.Context, c *db.Connection, itx *sql.Tx) (bool, error) {
		return c.Invoke(db.OptContext(ctx), db.OptTx(itx)).Query(fmt.Sprintf("select * from %s", tableName)).Any()
	})(
		context.Background(),
		defaultDB(),
		tx,
		action,
	)
	assert.Nil(err)
	assert.True(didRun)
}
