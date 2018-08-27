package stats

import (
	"context"

	"github.com/blend/go-sdk/web"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// TagKeyEnvironment is the environment (web, dev, etc.)
	TagKeyEnvironment = "env"
	// TagKeySpanType defines the Span type (web, db, cache).
	TagKeySpanType = "span.type"
	// TagKeyServiceName defines the Service name for this Span.
	TagKeyServiceName = "service.name"
	// TagKeyResourceName defines the Resource name for the Span.
	TagKeyResourceName = "resource.name"

	// SpanTypeWeb is a span type.
	SpanTypeWeb = "web"
	// SpanTypeDB is a span type.
	SpanTypeDB = "db"
	// SpanTypeSQL is a span type.
	SpanTypeSQL = "sql"
	// SpanTypeCache is a span type.
	SpanTypeCache = "cache"
	// SpanTypeRPC is a span type.
	SpanTypeRPC = "rpc"

	// TracingOperationHTTPRequest is the http request tracing operation name.
	TracingOperationHTTPRequest = "http.request"
	// TracingOperationDBPing is the db ping tracing operation.
	TracingOperationDBPing = "db.ping"
	// TracingOperationDBPrepare is the db prepare tracing operation.
	TracingOperationDBPrepare = "db.prepare"
	// TracingOperationDBQuery is the db query tracing operation.
	TracingOperationDBQuery = "db.query"
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

// GetTracingSpanFromContext returns a tracing span from a given context.
func GetTracingSpanFromContext(ctx context.Context, key string) opentracing.Span {
	if typed, ok := ctx.Value(key).(opentracing.Span); ok {
		return typed
	}
	return nil
}
