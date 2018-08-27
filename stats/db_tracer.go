package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/exception"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	_ db.Tracer = (*dbTracer)(nil)
)

// DBTracer returns a db tracer.
func DBTracer(tracer opentracing.Tracer) db.Tracer {
	return &dbTracer{tracer: tracer}
}

type dbTracer struct {
	tracer opentracing.Tracer
}

func (dbt dbTracer) Ping(ctx context.Context, conn *db.Connection) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeyResourceName, Value: conn.Config().GetDatabase()},
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeSQL},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := StartSpanFromContext(ctx, dbt.tracer, TracingOperationSQLPing, startOptions...)
	return dbTraceFinisher{span: span}
}

func (dbt dbTracer) Prepare(ctx context.Context, conn *db.Connection, statement string) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeyResourceName, Value: conn.Config().GetDatabase()},
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeSQL},
		opentracing.Tag{Key: "db.query", Value: statement},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := StartSpanFromContext(ctx, dbt.tracer, TracingOperationSQLPrepare, startOptions...)
	return dbTraceFinisher{span: span}
}

func (dbt dbTracer) Query(ctx context.Context, conn *db.Connection, inv *db.Invocation, statement string) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: TagKeyResourceName, Value: inv.Label()},
		opentracing.Tag{Key: TagKeySpanType, Value: SpanTypeSQL},
		opentracing.Tag{Key: TagKeyDBName, Value: conn.Config().GetDatabase()},
		opentracing.Tag{Key: TagKeyDBUser, Value: conn.Config().GetUsername()},
		opentracing.StartTime(inv.Start()),
	}
	span, _ := StartSpanFromContext(ctx, dbt.tracer, TracingOperationSQLQuery, startOptions...)
	return dbTraceFinisher{span: span}
}

type dbTraceFinisher struct {
	span opentracing.Span
}

func (dbtf dbTraceFinisher) Finish(err error) {
	if dbtf.span == nil {
		return
	}
	if err != nil {
		if typed := exception.As(err); typed != nil {
			dbtf.span.SetTag(TagKeyError, typed.Class())
			dbtf.span.SetTag(TagKeyErrorMessage, typed.Message())
			dbtf.span.SetTag(TagKeyErrorStack, typed.Stack().String())
		} else {
			dbtf.span.SetTag(TagKeyError, fmt.Sprintf("%v", err))
		}
	}
	dbtf.span.Finish()
}
