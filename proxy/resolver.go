package proxy

import (
	"net/http"
	"sync/atomic"
)

// Resolver is a function that takes a request and produces a destination `url.URL`.
type Resolver func(*http.Request, []*Upstream) (*Upstream, error)

// RoundRobinResolver returns a closure based resolver that rotates through upstreams uniformly.
func RoundRobinResolver(upstreams []*Upstream) Resolver {
	if len(upstreams) == 0 {
		return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
			return nil, nil
		}
	}

	if len(upstreams) == 1 {
		return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
			return upstreams[0], nil
		}
	}

	return manyRoundRobinResolver(upstreams)
}

func manyRoundRobinResolver(upstreams []*Upstream) Resolver {
	var current int32
	var total = int32(len(upstreams))
	return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
		index := atomic.LoadInt32(&current)
		upstream := upstreams[index]
		atomic.StoreInt32(&current, (index+1)%total)
		return upstream, nil
	}
}
