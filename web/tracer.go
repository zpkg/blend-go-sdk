package web

// Tracer is a type that listens to start and end events for a ctx.
type Tracer interface {
	Start(*Ctx)
	Finish(*Ctx)
}
