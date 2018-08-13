package web

import (
	"context"
	"time"
)

// Cancel injects the context for a given action with a cancel func.
func Cancel(action Action) Action {
	return func(ctx *Ctx) Result {
		ctx.ctx, ctx.cancel = context.WithCancel(context.Background())
		return action(ctx)
	}
}

// Timeout injects the context for a given action with a timeout context.
func Timeout(d time.Duration) Middleware {
	return func(action Action) Action {
		return func(ctx *Ctx) Result {
			ctx.ctx, ctx.cancel = context.WithTimeout(context.Background(), d)
			return action(ctx)
		}
	}
}

// View sets the context.DefaultResultProvider() equal to context.View().
func View(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.View()))
	}
}

// JSON sets the context.DefaultResultProvider() equal to context.JSON().
func JSON(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.JSON()))
	}
}

// XML sets the context.DefaultResultProvider() equal to context.XML().
func XML(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.XML()))
	}
}

// Text sets the context.DefaultResultProvider() equal to context.Text().
func Text(action Action) Action {
	return func(ctx *Ctx) Result {
		return action(ctx.WithDefaultResultProvider(ctx.Text()))
	}
}
