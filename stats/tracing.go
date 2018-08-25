package stats

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// StateKeySpan is the span state key.
	StateKeySpan = "span"
)

// TracingMiddleware returns a go-web middleware that creates a span for the routes action.
func TracingMiddleware(tracer opentracing.Tracer) web.Middleware {
	return func(action web.Action) web.Action {
		return func(ctx *web.Ctx) web.Result {
			// if the tracer is unset, just call the action ...
			if tracer == nil {
				return action(ctx)
			}

			// first try to extract the span from existing headers
			spanContext, err := tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(ctx.Request().Header))
			if err != nil {
				return ctx.DefaultResultProvider().BadRequest(err)
			}

			// open the span, inject headers in case we're opening a root span.
			var span opentracing.Span
			startOptions := []opentracing.StartSpanOption{
				opentracing.Tag{Key: "remote_addr", Value: logger.GetRemoteAddr(ctx.Request())},
				opentracing.Tag{Key: "host", Value: logger.GetHost(ctx.Request())},
				opentracing.Tag{Key: "user_agent", Value: logger.GetUserAgent(ctx.Request())},
				opentracing.StartTime(ctx.Start()),
			}
			if spanContext != nil {
				span = tracer.StartSpan(string(logger.HTTPRequest), append(startOptions, opentracing.ChildOf(spanContext))...)
			} else {
				span = tracer.StartSpan(string(logger.HTTPRequest), startOptions...)
				tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request().Header))
			}

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
