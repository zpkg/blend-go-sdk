package stats

import (
	"context"
	"time"

	"github.com/blend/go-sdk/db"
	opentracing "github.com/opentracing/opentracing-go"
)

// DBTracer returns a db tracer.
func DBTracer(tracer opentracing.Tracer) db.Tracer {
	return &dbTracer{tracer: tracer}
}

type dbTracer struct {
	tracer opentracing.Tracer
}

func (dbt dbTracer) PrepareStart(ctx context.Context, conn *db.Connection, statement string) {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeyResourceName, Value: conn.Config().GetDatabase()},
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeDB},
		opentracing.StartTime(time.Now().UTC()),
	}

	// start the span.
	span, spanCtx := StartSpanFromContext(ctx, dbt.tracer, TracingOperationDBPrepare, startOptions...)

}

func (dbt dbTracer) PrepareFinish(ctx context.Context, conn *db.Connection, statement string, err error) {

}

func (dbt dbTracer) InvocationStart(ctx context.Context, conn *db.Connection, inv *db.Invocation) {

}

func (dbt dbTracer) InvocationFinish(ctx context.Context, conn *db.Connection, inv *db.Invocation, statement string, err error) {

}
