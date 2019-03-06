package grpcutil

import (
	"net"

	"github.com/blend/go-sdk/async"
	"google.golang.org/grpc"
)

// NewGraceful returns a new graceful host for a grpc server.
func NewGraceful(listener net.Listener, server *grpc.Server) *Graceful {
	return &Graceful{
		Listener: listener,
		Server:   server,
	}
}

// Graceful is a shim for graceful hosting grpc servers.
type Graceful struct {
	latch    async.Latch
	Listener net.Listener
	Server   *grpc.Server
}

// Start starts the server.
func (gz *Graceful) Start() error {
	gz.latch.Starting()
	gz.latch.Started()
	return gz.Server.Serve(gz.Listener)
}

// Stop shuts the server down.
func (gz *Graceful) Stop() error {
	gz.latch.Stopping()
	gz.Server.GracefulStop()
	gz.latch.Stopped()
	return nil
}

// IsRunning returns if the server is running.
func (gz *Graceful) IsRunning() bool {
	return gz.latch.IsRunning()
}

// NotifyStarted returns the notify started signal.
func (gz *Graceful) NotifyStarted() <-chan struct{} {
	return gz.latch.NotifyStarted()
}

// NotifyStopped returns the notify stopped signal.
func (gz *Graceful) NotifyStopped() <-chan struct{} {
	return gz.latch.NotifyStopped()
}
