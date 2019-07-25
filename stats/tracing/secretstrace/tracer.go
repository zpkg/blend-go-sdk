package secretstrace

import (
	"github.com/blend/go-sdk/secrets"
	"github.com/blend/go-sdk/stats/tracing"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
	"time"
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) secrets.Tracer {
	return &secretsTracer{tracer: tracer}
}

type secretsTracer struct {
	tracer opentracing.Tracer
}

func (st secretsTracer) Start(req *http.Request) secrets.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeVault},
		opentracing.Tag{Key: tracing.TagSecretsOperation, Value: parseOperation(req)},
		opentracing.Tag{Key: tracing.TagSecretsMethod, Value: parseMethod(req)},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(req.Context(), st.tracer, tracing.OperationVaultAPI, startOptions...)
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	st.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	return secretsTraceFinisher{span: span}
}

type secretsTraceFinisher struct {
	span opentracing.Span
}

func (stf secretsTraceFinisher) Finish(req *http.Request, res *http.Response, err error) {
	if stf.span == nil {
		return
	}
	tracing.SpanError(stf.span, err)
	stf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(res.StatusCode))
	stf.span.Finish()
}

func parseMethod(req *http.Request) string {
	if req != nil && req.URL != nil {
		return req.Method
	}
	return "UNKNOWN"
}


func parseOperation(req *http.Request) string {
	if req != nil && req.URL != nil {
		return req.URL.Path
	}
	return "UNKNOWN"
}
