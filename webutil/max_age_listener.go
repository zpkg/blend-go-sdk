/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/blend/go-sdk/expvar"
)

// NewMaxAgeListener returns a new max age listener.
func NewMaxAgeListener(listener net.Listener, headerCloseAfter, forceCloseAfter time.Duration) *MaxAgeListener {
	return &MaxAgeListener{
		Listener:          listener,
		HeaderCloseAfter:  headerCloseAfter,
		ForceCloseAfter:   forceCloseAfter,
		ConnsOpened:       new(expvar.Int),
		ConnsHeaderClosed: new(expvar.Int),
		ConnsForcedClosed: new(expvar.Int),
	}
}

// MaxAgeListener is an extesion of keep alive listener
// that returns connections with age tracking.
//
// It also has a `ConnState` delegate implementation
// to close connections that are older than a given duration
// when they next return to idle.
type MaxAgeListener struct {
	net.Listener

	HeaderCloseAfter time.Duration
	ForceCloseAfter  time.Duration

	ConnsOpened       *expvar.Int
	ConnsHeaderClosed *expvar.Int
	ConnsForcedClosed *expvar.Int
}

// ApplyServer applies a max age listener to a server.
func (mal *MaxAgeListener) ApplyServer(server *http.Server) {
	server.Handler = mal.WrapHandler(server.Handler)
	server.ConnContext = mal.ConnContext
	server.ConnState = mal.ConnState
}

// Accept implements net.Listener
func (mal *MaxAgeListener) Accept() (net.Conn, error) {
	nc, err := mal.Listener.Accept()
	if err != nil {
		return nil, err
	}
	mal.ConnsOpened.Add(1)
	return &MaxAgeConn{Conn: nc, Created: time.Now()}, nil
}

// ConnContext implements the conn context handler.
func (mal *MaxAgeListener) ConnContext(ctx context.Context, nc net.Conn) context.Context {
	return WithNetConn(ctx, nc)
}

// WrapHandler wraps a handler with connection close semantics if a connection is older
// than the max conn age.
func (mal *MaxAgeListener) WrapHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if mal.HeaderCloseAfter == 0 {
			handler.ServeHTTP(rw, req)
			return
		}
		if conn := GetNetConn(req.Context()); conn != nil {
			if typed, ok := conn.(Ager); ok {
				if typed.Age() > mal.HeaderCloseAfter {
					mal.ConnsHeaderClosed.Add(1)
					rw.Header().Set(HeaderConnection, ConnectionClose)
				}
			}
		}
		handler.ServeHTTP(rw, req)
		return
	})
}

// ConnState handles conn state, and will forceably close connections when they return to idle
// if they're older than the max age.
func (mal *MaxAgeListener) ConnState(nc net.Conn, cs http.ConnState) {
	if mal.ForceCloseAfter == 0 {
		return
	}
	// if the connection is newly "idle" and too old, close it
	switch cs {
	case http.StateIdle:
		if typed, ok := nc.(Ager); ok && typed != nil {
			if typed.Age() > mal.ForceCloseAfter {
				mal.ConnsForcedClosed.Add(1)
				_ = nc.Close()
			}
		}
	default:
		return
	}
}

type netConnKey struct{}

// WithNetConn returns a context with a max age conn as a value.
func WithNetConn(ctx context.Context, nc net.Conn) context.Context {
	return context.WithValue(ctx, netConnKey{}, nc)
}

// GetNetConn gets the connection off a context.
func GetNetConn(ctx context.Context) net.Conn {
	if value := ctx.Value(netConnKey{}); value != nil {
		if typed, ok := value.(net.Conn); ok {
			return typed
		}
	}
	return nil
}

// Ager is a type that provides an age.
type Ager interface {
	Age() time.Duration
}

// MaxAgeConn is a connection that establishes an age.
type MaxAgeConn struct {
	net.Conn
	Created time.Time
}

// Age returns the connection age.
func (mac *MaxAgeConn) Age() time.Duration {
	return time.Since(mac.Created)
}
