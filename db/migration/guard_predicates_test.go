package migration

import (
	"context"
	"database/sql"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"testing"
)

func TestGuardsReal(t *testing.T) {
	assert := assert.New(t)
	tx, err := db.Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	tName := "table_test_foo"
	cName := "constraint_foo"
	colName := "created_foo"
	iName := "index_foo"

	var didRun bool
	action := Actions(func(ctx context.Context, c *db.Connection, itx *sql.Tx, _ ...db.InvocationOption) error {
		didRun = true
		return nil
	})

	err = TableExists(tName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = TableNotExists(tName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("CREATE TABLE table_test_foo (id serial not null primary key, something varchar(32) not null)", tx)
	assert.Nil(err)

	didRun = false
	err = TableExists(tName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	didRun = false
	err = ConstraintExists(tName, cName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = ConstraintNotExists(tName, cName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("ALTER TABLE table_test_foo ADD CONSTRAINT constraint_foo UNIQUE (something)", tx)
	assert.Nil(err)

	didRun = false
	err = ConstraintExists(tName, cName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	didRun = false
	err = ColumnExists(tName, colName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = ColumnNotExists(tName, colName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("ALTER TABLE table_test_foo ADD COLUMN created_foo timestamp not null", tx)
	assert.Nil(err)

	didRun = false
	err = ColumnExists(tName, colName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	didRun = false
	err = IndexExists(tName, iName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = IndexNotExists(tName, iName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("CREATE INDEX index_foo ON table_test_foo(created_foo DESC)", tx)
	assert.Nil(err)

	didRun = false
	err = IndexExists(tName, iName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)
}

func TestGuardsRealSchema(t *testing.T) {
	assert := assert.New(t)
	tx, err := db.Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	sName := "schema_test_bar"
	colName := "created_bar"
	cName := "constraint_bar"
	iName := "index_bar"
	tName := "table_test_bar"

	err = db.Default().ExecInTx("CREATE SCHEMA schema_test_bar", tx)
	assert.Nil(err)

	var didRun bool
	action := Actions(func(ctx context.Context, c *db.Connection, itx *sql.Tx, _ ...db.InvocationOption) error {
		didRun = true
		return nil
	})

	err = TableExistsInSchema(sName, tName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = TableNotExistsInSchema(sName, tName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("CREATE TABLE schema_test_bar.table_test_bar (id serial not null primary key, something varchar(32) not null, created timestamp not null)", tx)
	assert.Nil(err)

	didRun = false
	err = TableExistsInSchema(sName, tName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	didRun = false
	err = ConstraintExistsInSchema(sName, tName, cName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = ConstraintNotExistsInSchema(sName, tName, cName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("ALTER TABLE schema_test_bar.table_test_bar ADD CONSTRAINT constraint_bar UNIQUE (something)", tx)
	assert.Nil(err)

	didRun = false
	err = ConstraintExistsInSchema(sName, tName, cName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	didRun = false
	err = ColumnExistsInSchema(sName, tName, colName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = ColumnNotExistsInSchema(sName, tName, colName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("ALTER TABLE schema_test_bar.table_test_bar ADD COLUMN created_bar timestamp not null", tx)
	assert.Nil(err)

	didRun = false
	err = ColumnExistsInSchema(sName, tName, colName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	didRun = false
	err = IndexExistsInSchema(sName, tName, iName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.False(didRun)

	err = IndexNotExistsInSchema(sName, tName, iName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)

	err = db.Default().ExecInTx("CREATE INDEX index_bar ON schema_test_bar.table_test_bar(created_bar DESC)", tx)
	assert.Nil(err)

	didRun = false
	err = IndexExistsInSchema(sName, tName, iName)(context.Background(), db.Default(), tx, action)
	assert.Nil(err)
	assert.True(didRun)
}