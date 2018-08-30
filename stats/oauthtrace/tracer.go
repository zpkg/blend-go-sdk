package oauthtrace

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/oauth"
	"github.com/blend/go-sdk/stats/tracing"
	"github.com/opentracing/opentracing-go"
)

// Tracer returns a oauth tracer.
func Tracer(tracer opentracing.Tracer) oauth.Tracer {
	return &oauthTracer{tracer: tracer}
}

type oauthTracer struct {
	tracer opentracing.Tracer
}

func (ot oauthTracer) Start(r *http.Request) oauth.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeHTTP},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(r.Context(), ot.tracer, "oauth", startOptions...)
	return oauthTraceFinisher{span: span}
}

type oauthTraceFinisher struct {
	span opentracing.Span
}

func (otf oauthTraceFinisher) Finish(r *http.Request, res *oauth.Result, err error) {
	if otf.span == nil {
		return
	}
	tracing.SpanError(otf.span, err)
	otf.span.Finish()
}
