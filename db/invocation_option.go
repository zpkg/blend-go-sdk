package db

import (
	"context"
	"database/sql"
)

// InvocationOption is an option for invocations.
type InvocationOption func(*Invocation)

// OptCachedPlanKey sets the CachedPlanKey on the invocation.
func OptCachedPlanKey(cacheKey string) InvocationOption {
	return func(i *Invocation) {
		i.CachedPlanKey = cacheKey
	}
}

// OptContext sets a context on an invocation.
func OptContext(ctx context.Context) InvocationOption {
	return func(i *Invocation) {
		if ctx == nil {
			i.Context = context.Background()
		} else {
			i.Context = ctx
		}
	}
}

// OptTx is an invocation option that sets the invocation transaction.
func OptTx(tx *sql.Tx) InvocationOption {
	return func(i *Invocation) {
		i.Tx = tx
	}
}
