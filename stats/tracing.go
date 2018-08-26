package stats

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// StateKeySpan is the span state key.
	StateKeySpan = "web-span"

	// TracingOperationHTTPRequest is the http request tracing operation name.
	TracingOperationHTTPRequest = string(logger.HTTPRequest)
)

// TracingMiddleware returns a go-web middleware that creates a span for the routes action.
func TracingMiddleware(tracer opentracing.Tracer) web.Middleware {
	return func(action web.Action) web.Action {
		return func(ctx *web.Ctx) web.Result {
			// if the tracer is unset, just call the action ...
			if tracer == nil {
				return action(ctx)
			}

			// open the span, inject headers in case we're opening a root span.
			startOptions := []opentracing.StartSpanOption{
				opentracing.Tag{Key: "remote_addr", Value: logger.GetRemoteAddr(ctx.Request())},
				opentracing.Tag{Key: "host", Value: logger.GetHost(ctx.Request())},
				opentracing.Tag{Key: "user_agent", Value: logger.GetUserAgent(ctx.Request())},
				opentracing.StartTime(ctx.Start()),
			}
			if ctx.Route() != nil {
				startOptions = append(startOptions, opentracing.Tag{Key: "route", Value: ctx.Route().String()})
			}
			spanContext, _ := tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(ctx.Request().Header))
			if spanContext != nil {
				startOptions = append(startOptions, opentracing.ChildOf(spanContext))
			}

			span, spanCtx := opentracing.StartSpanFromContext(ctx.Context(), TracingOperationHTTPRequest, startOptions...)

			ctx.WithContext(spanCtx)
			ctx.WithStateValue(StateKeySpan, span)
			// close the span on exit...
			defer span.Finish()
			// call the action ...
			return action(ctx)
		}
	}
}

// GetTracingSpan gets a tracing span from a request.
func GetTracingSpan(ctx *web.Ctx) (span opentracing.Span) {
	value := ctx.StateValue(StateKeySpan)
	if typed, ok := value.(opentracing.Span); ok {
		return typed
	}
	return nil
}
