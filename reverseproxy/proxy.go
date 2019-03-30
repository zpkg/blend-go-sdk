package reverseproxy

import (
	"fmt"
	"net/http"
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

const (
	// FlagProxyRequest is a logger flag.
	FlagProxyRequest = "proxy.request"
)

// New returns a new proxy.
func New() *Proxy {
	return &Proxy{}
}

// Proxy is a factory for a simple reverse proxy.
type Proxy struct {
	UpstreamHeaders http.Header
	Log             logger.Log
	Upstreams       []*Upstream
	Resolver        Resolver
}

// WithUpstreamHeader adds a single upstream header.
func (p *Proxy) WithUpstreamHeader(key, value string) *Proxy {
	if p.UpstreamHeaders == nil {
		p.UpstreamHeaders = http.Header{}
	}
	p.UpstreamHeaders.Set(key, value)
	return p
}

// WithUpstream adds an upstream by URL.
func (p *Proxy) WithUpstream(upstream *Upstream) *Proxy {
	p.Upstreams = append(p.Upstreams, upstream)
	return p
}

// ServeHTTP is the http entrypoint.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			if p.Log != nil {
				p.Log.Fatalf("%v", r)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", r)
			}
		}
	}()

	// set the default resolver if unset.
	if p.Resolver == nil {
		p.Resolver = RoundRobinResolver(p.Upstreams)
	}

	upstream, err := p.Resolver(req, p.Upstreams)

	if err != nil {
		logger.MaybeError(p.Log, err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	if upstream == nil {
		rw.WriteHeader(http.StatusBadGateway)
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
	for key, values := range p.UpstreamHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	upstream.ServeHTTP(rw, req)
}
