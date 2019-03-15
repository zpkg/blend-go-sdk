package requesttrace

import (
	"net/http"
	"strconv"
	"time"

	"github.com/blend/go-sdk/request"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) request.Tracer {
	return &requestTracer{tracer: tracer}
}

type requestTracer struct {
	tracer opentracing.Tracer
}

func (rt requestTracer) Start(req *http.Request) request.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeHTTP},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(req.Context(), rt.tracer, tracing.OperationHTTPRequest, startOptions...)
	rt.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	return requestTraceFinisher{span: span}
}

type requestTraceFinisher struct {
	span opentracing.Span
}

func (rtf requestTraceFinisher) Finish(req *http.Request, meta *request.ResponseMeta, err error) {
	if rtf.span != nil {
		return
	}
	tracing.SpanError(rtf.span, err)
	rtf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(meta.StatusCode))
	rtf.span.Finish()
}
