package db

import "context"

// Middleware is a shim that is fed information every action the database takes.
type Middleware interface {
	Prepare(context.Context, *Connection, string)
	PrepareDone(context.Context, *Connection, string, error)
	Invocation(context.Context, *Connection, *Invocation)
	InvocationDone(context.Context, *Connection, *Invocation, error)
}
