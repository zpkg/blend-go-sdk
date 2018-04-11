package proxy

import (
	"net/http"
	"sync"
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
	l := sync.Mutex{}
	index := 0
	total := len(upstreams)

	return func(req *http.Request, upstreams []*Upstream) (*Upstream, error) {
		l.Lock()
		upstream := upstreams[index]
		index = (index + 1) % total
		l.Unlock()
		return upstream, nil
	}
}
