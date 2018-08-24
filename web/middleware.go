package web

import (
	"context"
	"time"
)

// WithCancel injects the context for a given action with a cancel func.
// It allows you to cancel the request.
func WithCancel(action Action) Action {
	return func(ctx *Ctx) Result {
		ctx.ctx, ctx.cancel = context.WithCancel(ctx.Context())
		return action(ctx)
	}
}

// WithTimeout injects the context for a given action with a timeout context.
func WithTimeout(d time.Duration) Middleware {
	return func(action Action) Action {
		return func(ctx *Ctx) Result {
			ctx.ctx, ctx.cancel = context.WithTimeout(ctx.Context(), d)
			return action(ctx)
		}
	}
}

// ViewProviderAsDefault sets the context.DefaultResultProvider() equal to context.View().
func ViewProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.View()))
	}
}

// JSONProviderAsDefault sets the context.DefaultResultProvider() equal to context.JSON().
func JSONProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.JSON()))
	}
}

// XMLProviderAsDefault sets the context.DefaultResultProvider() equal to context.XML().
func XMLProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.XML()))
	}
}

// TextProviderAsDefault sets the context.DefaultResultProvider() equal to context.Text().
func TextProviderAsDefault(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.Text()))
	}
}
