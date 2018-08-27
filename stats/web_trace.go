package stats

import (
	"strconv"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// StateKeySpan is the span state key.
	StateKeySpan = "web-span"

	// TracingOperationHTTPRequest is the http request tracing operation name.
	TracingOperationHTTPRequest = "http.request"
)

// WebTracer returns a web tracer.
func WebTracer(tracer opentracing.Tracer) web.Tracer {
	return &webTracer{tracer: tracer}
}

type webTracer struct {
	tracer opentracing.Tracer
}

func (wt webTracer) Start(ctx *web.Ctx) {
	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: "span.type", Value: "http"},
		opentracing.Tag{Key: "http.method", Value: ctx.Request().Method},
		opentracing.Tag{Key: "http.url", Value: ctx.Request().RequestURI},
		opentracing.Tag{Key: "http.remote_addr", Value: logger.GetRemoteAddr(ctx.Request())},
		opentracing.Tag{Key: "http.host", Value: logger.GetHost(ctx.Request())},
		opentracing.Tag{Key: "http.user_agent", Value: logger.GetUserAgent(ctx.Request())},
		opentracing.StartTime(ctx.Start()),
	}
	if ctx.Route() != nil {
		startOptions = append(startOptions, opentracing.Tag{Key: "http.route", Value: ctx.Route().String()})
	}

	// try to extract an incoming span context
	// this is typically done if we're a service being called in a chain from another (more ancestral)
	// span context.
	spanContext, _ := wt.tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(ctx.Request().Header))
	if spanContext != nil {
		startOptions = append(startOptions, opentracing.ChildOf(spanContext))
	}
	// start the span.
	span, spanCtx := StartSpanFromContext(ctx.Context(), wt.tracer, TracingOperationHTTPRequest, startOptions...)

	// inject the new context
	ctx.WithContext(spanCtx)
	// also store the span in the request state
	ctx.WithStateValue(StateKeySpan, span)
}

func (wt webTracer) Finish(ctx *web.Ctx) {
	span := GetTracingSpanFromCtx(ctx)
	if span == nil {
		return
	}
	span.SetTag("http.status_code", strconv.Itoa(ctx.Response().StatusCode()))
	span.Finish()
}
