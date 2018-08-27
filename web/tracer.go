package web

// Tracer is a type that listens to start and end events for a ctx.
type Tracer interface {
	Start(*Ctx) TraceFinisher
}

// TraceFinisher finishes a trace.
type TraceFinisher interface {
	Finish(*Ctx)
}
