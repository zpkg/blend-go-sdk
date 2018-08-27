package stats

import (
	"context"

	"github.com/blend/go-sdk/web"
	opentracing "github.com/opentracing/opentracing-go"
)

// StartSpanFromContext creates a new span from a given context.
// It is required because opentracing relies on global state.
func StartSpanFromContext(ctx context.Context, tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := tracer.StartSpan(operationName, opts...)
	return span, opentracing.ContextWithSpan(ctx, span)
}

// GetTracingSpanFromCtx gets a tracing span from a web ctx.
func GetTracingSpanFromCtx(ctx *web.Ctx) (span opentracing.Span) {
	value := ctx.StateValue(StateKeySpan)
	if typed, ok := value.(opentracing.Span); ok {
		return typed
	}
	return nil
}
