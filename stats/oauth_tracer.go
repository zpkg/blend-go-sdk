package stats

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/oauth"
	"github.com/opentracing/opentracing-go"
)

// OAuthTracer returns a oauth tracer.
func OAuthTracer(tracer opentracing.Tracer) oauth.Tracer {
	return &oauthTracer{tracer: tracer}
}

type oauthTracer struct {
	tracer opentracing.Tracer
}

func (ot oauthTracer) Start(r *http.Request) oauth.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeHTTP},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := StartSpanFromContext(r.Context(), ot.tracer, "oauth", startOptions...)
	return oauthTraceFinisher{span: span}
}

type oauthTraceFinisher struct {
	span opentracing.Span
}

func (otf oauthTraceFinisher) Finish(r *http.Request, res *oauth.Result, err error) {
	if otf.span == nil {
		return
	}
	SpanError(otf.span, err)
	otf.span.Finish()
}
