package migration

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

// NewGroup creates a new GroupedActions from a given list of actionable.
func NewGroup(options ...GroupOption) *GroupedActions {
	g := &GroupedActions{}
	for _, o := range options {
		o(g)
	}
	return g
}

// NewWithActions is a helper for "migrations.NewGroup(migrations.OptActions(actions ...migration.Actionable))"
func NewWithActions(actions ...Actionable) *GroupedActions {
	return NewGroup(OptActions(actions...))
}

// GroupedActions is an atomic series of migration actions.
// It uses transactions to apply these actions as an atomic unit.
type GroupedActions struct{
	Actions []Actionable
	NoTransaction bool
}

// Action runs the groups actions within a transaction.
func (ga *GroupedActions) Action(ctx context.Context, c *db.Connection) (err error) {
	var tx *sql.Tx = nil
	if !ga.NoTransaction {
		tx, err = c.Begin()
		if err != nil {
			return
		}
	}

	// commit or rollback the transaction if this is a transactional operation
	if tx != nil {
		defer func() {
			if err != nil {
				if txErr := tx.Rollback(); txErr != nil {
					err = ex.Nest(err, txErr)
				}
			} else {
				if txErr := tx.Commit(); txErr != nil {
					err = ex.Nest(err, txErr)
				}
			}
		}()
	}

	for _, a := range ga.Actions {
		err = a.Action(ctx, c, tx)
		if err != nil {
			return
		}
	}

	return
}
