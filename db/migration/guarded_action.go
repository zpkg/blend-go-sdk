package migration

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/db"
)

// Step returns a new guarded actionable.
func Step(guard GuardFunc, action Action) *GuardedAction {
	return &GuardedAction{
		Guard: guard,
		Body:  action,
	}
}

// GuardedAction is a guarded actionable.
type GuardedAction struct {
	Guard GuardFunc
	Body  Action
}

// Action runs the body if the provided guard passes.
func (ga GuardedAction) Action(ctx context.Context, c *db.Connection, tx *sql.Tx) error {
	return ga.Guard(ctx, c, tx, ga.Body)
}
