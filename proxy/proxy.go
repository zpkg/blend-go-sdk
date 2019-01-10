package proxy

import (
	"fmt"
	"net/http"
	"os"

	"github.com/blend/go-sdk/logger"
)

const (
	// FlagProxyRequest is a logger flag.
	FlagProxyRequest logger.Flag = "proxy.request"
)

// New returns a new proxy.
func New() *Proxy {
	return &Proxy{}
}

// Proxy is a factory for a simple reverse proxy.
type Proxy struct {
	log       logger.Log
	upstreams []*Upstream
	resolver  Resolver
}

// WithLogger sets a property and returns the proxy reference.
func (p *Proxy) WithLogger(log logger.Log) *Proxy {
	p.log = log
	return p
}

// WithUpstream adds an upstream by URL.
func (p *Proxy) WithUpstream(upstream *Upstream) *Proxy {
	p.upstreams = append(p.upstreams, upstream)
	return p
}

// WithResolver sets a property and returns the proxy reference.
func (p *Proxy) WithResolver(resolver Resolver) *Proxy {
	p.resolver = resolver
	return p
}

// ServeHTTP is the http entrypoint.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			if p.log != nil {
				p.log.Fatalf("%v", r)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", r)
			}
		}
	}()

	// set the default resolver if unset.
	if p.resolver == nil {
		p.resolver = RoundRobinResolver(p.upstreams)
	}

	upstream, err := p.resolver(req, p.upstreams)

	if err != nil {
		p.log.Error(err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	if upstream == nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	upstream.ServeHTTP(rw, req)
}
