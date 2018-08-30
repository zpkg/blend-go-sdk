package crontrace

import (
	"context"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// Tracer returns a opentracing cron tracer.
func Tracer(t opentracing.Tracer) cron.Tracer {
	return &tracer{tracer: t}
}

type tracer struct {
	tracer opentracing.Tracer
}

func (t tracer) Start(ctx context.Context, task cron.Task) cron.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeJob},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: task.Name()},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, t.tracer, tracing.OperationJob, startOptions...)
	return &traceFinisher{span: span}
}

type traceFinisher struct {
	span opentracing.Span
}

func (tf traceFinisher) Finish(ctx context.Context, t cron.Task, err error) {
	if tf.span == nil {
		return
	}
	tracing.SpanError(tf.span, err)
	tf.span.Finish()
}
