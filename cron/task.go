package cron

import (
	"context"
)

// --------------------------------------------------------------------------------
// interfaces
// --------------------------------------------------------------------------------

// Task is an interface that structs can satisfy to allow them to be run as tasks.
type Task interface {
	Name() string
	Execute(ctx context.Context) error
}
