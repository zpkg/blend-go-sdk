/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net"
	"net/http"
	"time"
)

// MaxAgeListener is an extesion of keep alive listener
// that returns connections with age tracking.
//
// It also has a `ConnState` delegate implementation
// to close connections that are older than a given duration
// when they next return to idle.
type MaxAgeListener struct {
	net.Listener
	MaxConnAge time.Duration
}

// Accept implements net.Listener
func (mal MaxAgeListener) Accept() (net.Conn, error) {
	nc, err := mal.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return MaxAgeConn{Conn: nc, Created: time.Now()}, nil
}

// ConnState handles conn state.
func (mal MaxAgeListener) ConnState(nc net.Conn, cs http.ConnState) {
	if mal.MaxConnAge == 0 {
		return
	}
	// if the connection is newly "idle" and too old
	switch cs {
	case http.StateIdle:
		if typed, ok := nc.(MaxAger); ok {
			if typed.Age() > mal.MaxConnAge {
				_ = nc.Close()
			}
		}
	default:
		return
	}
}

// MaxAger is a type that provides an age.
type MaxAger interface {
	Age() time.Duration
}

// MaxAgeConn is a connection that establishes an age.
type MaxAgeConn struct {
	net.Conn
	Created time.Time
}

// Age returns the connection age.
func (mac MaxAgeConn) Age() time.Duration {
	return time.Since(mac.Created)
}
