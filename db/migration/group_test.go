package migration

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/util"
)

func TestGroupAbortOnError(t *testing.T) {
	assert := assert.New(t)

	// test if a migration group aborts the step list on an error from a child step.

	var didRun bool
	var txWasSet bool
	group := NewGroup(
		NewStep(AlwaysRun(), Actions(func(c *db.Connection, tx *sql.Tx) error {
			txWasSet = tx != nil
			return nil
		})),
		NewStep(AlwaysRun(), Actions(func(c *db.Connection, tx *sql.Tx) error {
			return fmt.Errorf("only a test")
		})),
		NewStep(AlwaysRun(), Actions(func(c *db.Connection, tx *sql.Tx) error {
			didRun = true
			return nil
		})),
	).WithAbortOnError(true)

	assert.True(group.AbortOnError())
	assert.NotNil(group.Apply(db.Default()))
	assert.False(didRun)
	assert.True(txWasSet)

	group.WithTransactionBound(false)
	assert.True(group.AbortOnError())
	assert.NotNil(group.Apply(db.Default()))
	assert.False(didRun)
}

func TestGroupTransactionBound(t *testing.T) {
	assert := assert.New(t)

	// test if a migration group shares a transaction for child steps.
	// the old behavior is each step got it's own transaction.

	tableName := util.String.RandomLetters(24)
	group := NewGroup(
		NewStep(AlwaysRun(), Statements(fmt.Sprintf(`CREATE TABLE %s (name varchar(255))`, tableName))),
		NewStep(AlwaysRun(), Actions(func(c *db.Connection, tx *sql.Tx) error {
			_, err := c.QueryInTx(fmt.Sprintf("select 1 from %s", tableName), tx).Any()
			return err
		})),
		NewStep(AlwaysRun(), Statements(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName))),
	).WithAbortOnError(false)

	group.WithTransactionBound(true)

	assert.False(group.AbortOnError())
	assert.True(group.TransactionBound())
	assert.Nil(group.Apply(db.Default()))
}

func TestGroupTransactionUnbound(t *testing.T) {
	assert := assert.New(t)

	// test if we skip passing transactions to transaction unbound steps.

	var didRun bool
	var unboundTXWasSet, boundTxWasSet bool
	group := NewGroup(
		NewStep(AlwaysRun(), Actions(func(c *db.Connection, tx *sql.Tx) error {
			unboundTXWasSet = tx != nil
			return nil
		})).WithTransactionBound(false),
		NewStep(AlwaysRun(), Actions(func(c *db.Connection, tx *sql.Tx) error {
			didRun = true
			boundTxWasSet = tx != nil
			return nil
		})),
	).WithAbortOnError(true)

	assert.True(group.AbortOnError())
	assert.True(group.TransactionBound())
	assert.Nil(group.Apply(db.Default()))
	assert.False(unboundTXWasSet)
	assert.True(didRun)
	assert.True(boundTxWasSet)

	group.WithTransactionBound(false)
	assert.True(group.AbortOnError())
	assert.Nil(group.Apply(db.Default()))
	assert.False(unboundTXWasSet)
	assert.True(didRun)
	assert.False(boundTxWasSet)
}

func TestGroupRollbackOnComplete(t *testing.T) {
	// test if we roll back changes on success (useful for testing).
}
