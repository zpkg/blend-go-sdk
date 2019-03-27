package pg

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
)

func TestGuard(t *testing.T) {
	assert := assert.New(t)
	tx, err := db.Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	tableName := randomName()
	err = createTestTable(tableName, tx)
	assert.Nil(err)

	err = insertTestValue(tableName, 4, "test", tx)
	assert.Nil(err)

	var didRun bool
	action := migration.Actions(func(ctx context.Context, c *db.Connection, itx *sql.Tx) error {
		didRun = true
		return nil
	})

	err = migration.Guard("test", func(c *db.Connection, itx *sql.Tx) (bool, error) {
		return c.Invoke(context.Background(), db.OptTx(itx)).Query(fmt.Sprintf("select * from %s", tableName)).Any()
	})(
		context.Background(),
		db.Default(),
		tx,
		action,
	)
	assert.Nil(err)
	assert.True(didRun)
}
