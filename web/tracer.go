package web

import "net/http"

// Tracer is a type that traces complete requests.
type Tracer interface {
	Start(*Ctx) TraceFinisher
}

// TraceFinisher is a finisher for a trace.
type TraceFinisher interface {
	Finish(*Ctx, error)
}

// RouteTracer is a type that can listen for route lookup traces.
type RouteTracer interface {
	StartRoute(*http.Request) RouteTraceFinisher
}

// RouteTraceFinisher is a finisher for route lookup traces.
type RouteTraceFinisher interface {
	Finish(*http.Request, string)
}

// ViewTracer is a type that can listen for view rendering traces.
type ViewTracer interface {
	StartView(*Ctx, *ViewResult) ViewTraceFinisher
}

// ViewTraceFinisher is a finisher for view traces.
type ViewTraceFinisher interface {
	Finish(*Ctx, *ViewResult, error)
}
