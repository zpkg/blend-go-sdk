/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package dbtrace

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/zpkg/blend-go-sdk/db"
	"github.com/zpkg/blend-go-sdk/tracing"
)

var (
	_ db.Tracer = (*dbTracer)(nil)
)

// Tracer returns a db tracer.
func Tracer(tracer opentracing.Tracer) db.Tracer {
	return &dbTracer{tracer: tracer}
}

type dbTracer struct {
	tracer opentracing.Tracer
}

func (dbt dbTracer) Prepare(ctx context.Context, cfg db.Config, statement string) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeSQL},
		opentracing.Tag{Key: tracing.TagKeyDBName, Value: cfg.DatabaseOrDefault()},
		opentracing.Tag{Key: tracing.TagKeyDBUser, Value: cfg.Username},
		opentracing.Tag{Key: TagKeyQuery, Value: statement},
		opentracing.Tag{Key: TagKeySQLCommand, Value: statement},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, dbt.tracer, tracing.OperationSQLPrepare, startOptions...)
	return dbTraceFinisher{span: span}
}

func (dbt dbTracer) Query(ctx context.Context, cfg db.Config, label, statement string) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: label},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeSQL},
		opentracing.Tag{Key: tracing.TagKeyDBName, Value: cfg.DatabaseOrDefault()},
		opentracing.Tag{Key: tracing.TagKeyDBUser, Value: cfg.Username},
		opentracing.Tag{Key: TagKeyQuery, Value: statement},
		opentracing.Tag{Key: TagKeySQLCommand, Value: statement},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, dbt.tracer, tracing.OperationSQLQuery, startOptions...)
	return dbTraceFinisher{span: span}
}

type dbTraceFinisher struct {
	span opentracing.Span
}

func (dbtf dbTraceFinisher) FinishPrepare(ctx context.Context, err error) {
	if dbtf.span == nil {
		return
	}
	if err == driver.ErrSkip {
		return
	}
	tracing.SpanError(dbtf.span, err)
	dbtf.span.Finish()
}

func (dbtf dbTraceFinisher) FinishQuery(ctx context.Context, res sql.Result, err error) {
	if dbtf.span == nil {
		return
	}
	if err == driver.ErrSkip {
		return
	}
	if res != nil {
		affected, _ := res.RowsAffected()
		dbtf.span.SetTag(tracing.TagKeyDBRowsAffected, affected)
	}
	tracing.SpanError(dbtf.span, err)
	dbtf.span.Finish()
}
