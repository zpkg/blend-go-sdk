package grpcutil

import "github.com/blend/go-sdk/graceful"

// Validate the interface is satisfied.
var (
	_ (graceful.Graceful) = (*Graceful)(nil)
)
