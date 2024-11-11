/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zpkg/blend-go-sdk/ex"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/webutil"
)

const (
	// FlagProxyRequest is a logger flag.
	FlagProxyRequest = "proxy.request"
)

// NewProxy returns a new proxy.
func NewProxy(opts ...ProxyOption) (*Proxy, error) {
	var err error
	p := Proxy{
		Headers: http.Header{},
	}
	for _, opt := range opts {
		if err = opt(&p); err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// Proxy is a factory for a simple reverse proxy.
type Proxy struct {
	Headers          http.Header
	Log              logger.Log
	Upstreams        []*Upstream
	Resolver         Resolver
	Tracer           webutil.HTTPTracer
	TransformRequest TransformRequest
	Timeout          time.Duration
}

// ServeHTTP is the http entrypoint.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var err error
	var tf webutil.HTTPTraceFinisher
	srw := webutil.NewStatusResponseWriter(rw)

	defer func() {
		// NOTE: This uses the outer scope's `err` by design. This way updates
		//       to `err` will be reflected on (deferred) exit.
		r := recover()

		// see: https://golang.org/pkg/net/http/#ErrAbortHandler
		if r != nil && r != http.ErrAbortHandler {
			// Wrap the error with the reason for the panic.
			err = ex.Nest(err, ex.New(r))
		}
		if tf != nil {
			tf.Finish(srw.StatusCode(), err)
		}
		if err != nil {
			if p.Log != nil {
				p.Log.Fatalf("%v", r)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", r)
			}
		}
	}()

	if p.Tracer != nil {
		tf, req = p.Tracer.Start(req)
	}

	// set the default resolver if unset.
	if p.Resolver == nil {
		p.Resolver = RoundRobinResolver(p.Upstreams)
	}

	upstream, err := p.Resolver(req, p.Upstreams)
	if err != nil {
		logger.MaybeError(p.Log, err)
		srw.WriteHeader(http.StatusBadGateway)
		return
	}

	if upstream == nil {
		srw.WriteHeader(http.StatusBadGateway)
		return
	}

	// Add extra forwarded headers.
	// these are required for a majority of services to function correctly behind
	// a reverse proxy.
	// They are "Add" vs. "Set" in case there are existing values.
	if port := webutil.GetPort(req); port != "" {
		req.Header.Add("X-Forwarded-Port", port)
	}
	if proto := webutil.GetProto(req); proto != "" {
		req.Header.Add("X-Forwarded-Proto", proto)
	}
	// add upstream headers.
	for key, values := range p.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if p.TransformRequest != nil {
		p.TransformRequest(req)
	}
	upstream.ServeHTTP(srw, req)
}
