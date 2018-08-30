package cron

import "context"

// Tracer is a trace handler.
type Tracer interface {
	Start(context.Context, Task) (context.Context, TraceFinisher)
}

// TraceFinisher is a finisher for traces.
type TraceFinisher interface {
	Finish(context.Context, Task, error)
}
