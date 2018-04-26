package migration

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
)

func TestGroup(t *testing.T) {
	assert := assert.New(t)

	g := NewGroup()
	assert.False(g.TransactionBound())
	assert.True(g.AbortOnError())
	assert.True(g.IsRoot())
	assert.False(g.RollbackOnComplete())
	assert.Empty(g.Label())
	assert.Nil(g.Collector())

	g.With(NewStep(AlwaysRun(), NoOp))
	assert.Len(g.migrations, 1)

	g.WithLabel("test")
	assert.Equal("test", g.Label())

	g.WithParent(NewGroup().WithLabel("parent"))
	assert.NotNil(g.Parent())
	assert.False(g.IsRoot())

	g.WithAbortOnError(false)
	assert.False(g.AbortOnError())

	g.WithCollector(&Collector{})
	assert.NotNil(g.Collector())

	assert.Nil(g.Collector().output)
	g.WithLogger(logger.None())
	assert.NotNil(g.Collector().output)

	g.WithUseTransaction(false)
	assert.False(g.UseTransaction())
}

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

	group.WithUseTransaction(false)
	assert.True(group.AbortOnError())
	assert.NotNil(group.Apply(db.Default()))
	assert.False(didRun)
}

func TestGroupUseTransaction(t *testing.T) {
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

	group.WithUseTransaction(true)

	assert.False(group.AbortOnError())
	assert.True(group.UseTransaction())
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
	assert.False(group.TransactionBound())
	assert.True(group.UseTransaction())
	assert.Nil(group.Apply(db.Default()))
	assert.False(unboundTXWasSet)
	assert.True(didRun)
	assert.True(boundTxWasSet)

	group.WithUseTransaction(false)
	assert.True(group.AbortOnError())
	assert.Nil(group.Apply(db.Default()))
	assert.False(unboundTXWasSet)
	assert.True(didRun)
	assert.False(boundTxWasSet)
}

func TestGroupRollbackOnComplete(t *testing.T) {
	// test if we roll back changes on success (useful for testing).
}
