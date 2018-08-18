package migration

import (
	"database/sql"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/exception"
)

// NewGroup creates a new migration group.
func NewGroup(steps ...*Step) *Group {
	r := &Group{
		useTransaction: true,
		abortOnError:   true,
	}
	return r.With(migrations...)
}

// Group is an atomic series of migrations.
// It uses transactions to apply a set of sub-migrations as a unit.
type Group struct {
	label  string
	parent *Suite
	steps  []*Step
}

// WithSteps adds steps to the group.
func (g *Group) WithSteps(steps ...*Step) {
	for _, s := range steps {
		g.steps = append(g.steps, s.WithParent(g))
	}
	return g
}

// WithLabel sets the migration label.
func (g *Group) WithLabel(value string) Migration {
	g.label = value
	return g
}

// Label returns a label for the runner.
func (g *Group) Label() string {
	return g.label
}

// WithParent sets the runner's parent.
func (g *Group) WithParent(parent *Suite) *Group {
	g.parent = parent
	return g
}

// Parent returns the runner's parent.
func (g *Group) Parent() *Group {
	return g.parent
}

// Apply wraps the action in a transaction and commits it if there were no errors, rolling back if there were.
func (g *Group) Apply(c *db.Connection, txs ...*sql.Tx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}

		if g.IsRoot() && g.collector != nil {
			g.collector.WriteStats()
		}
	}()

	// if the migration
	var tx *sql.Tx
	if g.useTransaction {
		tx, err = c.Begin()
		if err != nil {
			return
		}

		defer func() {
			if err != nil || g.rollbackOnComplete {
				err = exception.Nest(err, exception.New(tx.Rollback()))
			} else if err == nil {
				err = exception.New(tx.Commit())
			}
		}()
	}

	for _, m := range g.migrations {
		// if the migration is a group or something else that manages it's own transactions ...
		if m.TransactionBound() {
			err = m.WithCollector(g.collector).Apply(c, tx)
		} else {
			err = m.WithCollector(g.collector).Apply(c)
		}
		if err != nil && g.abortOnError {
			return
		}
		continue
	}

	return
}
