package r2

import "net/http"

// Tracer is a tracer for requests.
type Tracer interface {
	Start(*http.Request) TraceFinisher
}

// TraceFinisher is a finisher for traces.
type TraceFinisher interface {
	Finish(*http.Request, *http.Response, error)
}
