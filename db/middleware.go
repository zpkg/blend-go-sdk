package db

import (
	"context"
	"time"
)

// Middleware is a shim that is fed information every action the database takes.
type Middleware interface {
	Prepare(ctx context.Context, conn *Connection, statement string)
	PrepareDone(ctx context.Context, conn *Connection, statement string, elapsed time.Duration, err error)
	Invocation(ctx context.Context, invocation *Invocation, statement string)
	InvocationDone(ctx context.Context, invocation *Invocation, statement string, elapsed time.Duration, err error)
}
