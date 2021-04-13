/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"sync"
	"time"
)

var (
	_ http.RoundTripper = (*MaxAgeTransport)(nil)
)

// MaxAgeTransport is a wrapper for `http.Transport` that
// implements keep alive max connection age.
type MaxAgeTransport struct {
	http.Transport
	MaxConnAge time.Duration

	closeWorkerOnce sync.Once
	closeWorkerStop chan struct{}
}

// RoundTrip implements http.RoundTripper.
func (mat *MaxAgeTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	mat.closeWorkerOnce.Do(mat.startCloseWorker)
	return mat.Transport.RoundTrip(req)
}

func (mat *MaxAgeTransport) startCloseWorker() {
	if mat.MaxConnAge <= 0 {
		panic("invalid max connection age")
	}
	go mat.closeWorker()
}

func (mat *MaxAgeTransport) closeWorker() {
	mat.closeWorkerStop = make(chan struct{})
	ticker := time.NewTicker(mat.MaxConnAge)
	for {
		select {
		case <-ticker.C:
			mat.Transport.CloseIdleConnections()
		case <-mat.closeWorkerStop:
			return
		}
	}
}

// Close shuts down the closer worker if it's started.
func (mat *MaxAgeTransport) Close() error {
	if mat.closeWorkerStop != nil {
		close(mat.closeWorkerStop)
	}
	return nil
}
