package migration

import (
	"database/sql"

	"github.com/blend/go-sdk/db"
)

// Guard is a control for migration steps.
type Guard func(s *Step, c *db.Connection, tx *sql.Tx) error

// Invocable is a thing that can be invoked.
type Invocable interface {
	Invoke(c *db.Connection, tx *sql.Tx) error
}

// InvocableAction is a function that can be run during a migration step.
type InvocableAction func(c *db.Connection, tx *sql.Tx) error
