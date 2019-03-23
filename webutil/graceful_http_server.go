package webutil

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
)

// NewGracefulHTTPServer returns a new graceful http server wrapper.
func NewGracefulHTTPServer(server *http.Server, options ...GracefulHTTPServerOption) *GracefulHTTPServer {
	return &GracefulHTTPServer{
		Latch:  async.NewLatch(),
		Server: server,
	}
}

// GracefulHTTPServerOption is an option for the graceful http server.
type GracefulHTTPServerOption func(*GracefulHTTPServer)

// OptGracefulHTTPServerShutdownGracePeriod sets the shutdown grace period.
func OptGracefulHTTPServerShutdownGracePeriod(d time.Duration) GracefulHTTPServerOption {
	return func(g *GracefulHTTPServer) { g.ShutdownGracePeriod = d }
}

// OptGracefulHTTPServerListener sets the server listener.
func OptGracefulHTTPServerListener(listener net.Listener) GracefulHTTPServerOption {
	return func(g *GracefulHTTPServer) { g.Listener = listener }
}

// GracefulHTTPServer is a wrapper for an http server that implements the graceful interface.
type GracefulHTTPServer struct {
	*async.Latch
	Server              *http.Server
	ShutdownGracePeriod time.Duration
	Listener            net.Listener
}

// Start implements graceful.Graceful.Start.
// It is expected to block.
func (gs *GracefulHTTPServer) Start() (err error) {
	if !gs.CanStart() {
		err = exception.New(async.ErrCannotStart)
		return
	}
	gs.Starting()
	gs.Started()
	defer gs.Stopped()

	var shutdownErr error
	if gs.Listener != nil {
		shutdownErr = gs.Server.Serve(gs.Listener)
	} else {
		shutdownErr = gs.Server.ListenAndServe()
	}
	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		err = exception.New(shutdownErr)
	}
	return
}

// Stop implements graceful.Graceful.Stop.
func (gs *GracefulHTTPServer) Stop() error {
	if !gs.CanStop() {
		return exception.New(async.ErrCannotStop)
	}
	gs.Stopping()
	gs.Server.SetKeepAlivesEnabled(false)
	ctx := context.Background()
	if gs.ShutdownGracePeriod > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, gs.ShutdownGracePeriod)
		defer cancel()
	}
	return exception.New(gs.Server.Shutdown(ctx))
}
