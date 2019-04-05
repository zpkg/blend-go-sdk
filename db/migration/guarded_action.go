package migration

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/db"
)

// Step returns a new guarded actionable.
func Step(guard GuardFunc, action Action, options ...db.InvocationOption) *GuardedAction {
	return &GuardedAction{
		Guard:   guard,
		Body:    action,
		Options: options,
	}
}

// GuardedAction is a guarded actionable.
type GuardedAction struct {
	Guard   GuardFunc
	Body    Action
	Options []db.InvocationOption
}

// BodyWithOptions is the guarded action body with a given set of options.
func (ga GuardedAction) BodyWithOptions(ctx context.Context, c *db.Connection, tx *sql.Tx, options ...db.InvocationOption) error {
	return ga.Body(ctx, c, tx, append(options, ga.Options...)...)
}

// Action runs the body if the provided guard passes.
func (ga GuardedAction) Action(ctx context.Context, c *db.Connection, tx *sql.Tx, options ...db.InvocationOption) error {
	return ga.Guard(ctx, c, tx, ga.BodyWithOptions)
}
