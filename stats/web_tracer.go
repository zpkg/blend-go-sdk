package stats

import (
	"strconv"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// StateKeySpan is the span state key.
	StateKeySpan = "web-span"
)

var (
	_ web.Tracer = (*webTracer)(nil)
)

// WebTracer returns a web tracer.
func WebTracer(tracer opentracing.Tracer) web.Tracer {
	return &webTracer{tracer: tracer}
}

type webTracer struct {
	tracer opentracing.Tracer
}

func (wt webTracer) Start(ctx *web.Ctx) web.TraceFinisher {
	var resource string
	if ctx.Route() != nil {
		resource = ctx.Route().String()
	} else {
		resource = ctx.Request().URL.Path
	}

	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeyResourceName, Value: resource},
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeWeb},
		opentracing.Tag{Key: TagKeyHTTPMethod, Value: ctx.Request().Method},
		opentracing.Tag{Key: TagKeyHTTPURL, Value: ctx.Request().URL.Path},
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
	ctx.Request().WithContext(spanCtx)
	ctx.WithContext(spanCtx)

	// also store the span in the request state
	ctx.WithStateValue(StateKeySpan, span)
	return &webTraceFinisher{span: span}
}

type webTraceFinisher struct {
	span opentracing.Span
}

func (wtf webTraceFinisher) Finish(ctx *web.Ctx, err error) {
	if wtf.span == nil {
		return
	}
	SpanError(wtf.span, err)
	wtf.span.SetTag(TagKeyHTTPCode, strconv.Itoa(ctx.Response().StatusCode()))
	wtf.span.Finish()
}

func (wt webTracer) StartView(ctx *web.Ctx, vr *web.ViewResult) web.ViewTraceFinisher {
	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeyResourceName, Value: vr.ViewName},
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeWeb},
		opentracing.StartTime(time.Now().UTC()),
	}
	// start the span.
	span, _ := StartSpanFromContext(ctx.Context(), wt.tracer, TracingOperationHTTPRender, startOptions...)
	return &webViewTraceFinisher{span: span}
}

type webViewTraceFinisher struct {
	span opentracing.Span
}

func (wvtf webViewTraceFinisher) Finish(ctx *web.Ctx, vr *web.ViewResult, err error) {
	if wvtf.span == nil {
		return
	}
	SpanError(wvtf.span, err)
	wvtf.span.Finish()
}
