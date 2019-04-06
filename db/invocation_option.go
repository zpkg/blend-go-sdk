package db

import (
	"context"
	"database/sql"
	"time"
)

// InvocationOption is an option for invocations.
type InvocationOption func(*Invocation)

// OptCachedPlanKey sets the CachedPlanKey on the invocation.
func OptCachedPlanKey(cacheKey string) InvocationOption {
	return func(i *Invocation) {
		i.CachedPlanKey = cacheKey
	}
}

// OptInvocationStatementInterceptor sets the invocation statement interceptor.
func OptInvocationStatementInterceptor(interceptor StatementInterceptor) InvocationOption {
	return func(i *Invocation) { i.StatementInterceptor = interceptor }
}

// OptContext sets a context on an invocation.
func OptContext(ctx context.Context) InvocationOption {
	return func(i *Invocation) {
		i.Context = ctx
	}
}

// OptCancel sets the context cancel func..
func OptCancel(cancel context.CancelFunc) InvocationOption {
	return func(i *Invocation) {
		i.Cancel = cancel
	}
}

// OptTimeout sets a command timeout for the invocation.
func OptTimeout(d time.Duration) InvocationOption {
	return func(i *Invocation) {
		i.Context, i.Cancel = context.WithTimeout(i.Context, d)
	}
}

// OptTx is an invocation option that sets the invocation transaction.
func OptTx(tx *sql.Tx) InvocationOption {
	return func(i *Invocation) {
		i.Tx = tx
	}
}
