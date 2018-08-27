package db

import "context"

// Tracer is a shim that can be used to report stats for queries.
type Tracer interface {
	PrepareStart(context.Context, *Connection, string)
	PrepareFinish(context.Context, *Connection, string, error)

	InvocationStart(context.Context, *Connection, *Invocation)
	InvocationFinish(context.Context, *Connection, *Invocation, string, error)
}
