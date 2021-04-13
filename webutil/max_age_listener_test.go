/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_MaxAgeListener(t *testing.T) {
	its := assert.New(t)

	defaultMaxConnAge := 50 * time.Millisecond

	l, err := net.Listen("tcp", "127.0.0.1:0")
	its.Nil(err)

	mal := MaxAgeListener{
		Listener:   l,
		MaxConnAge: defaultMaxConnAge,
	}

	var serverCalls int32
	mockServer := &http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			atomic.AddInt32(&serverCalls, 1)
			rw.WriteHeader(http.StatusOK)
			fmt.Fprintf(rw, "OK!\n")
		}),
		ConnState: mal.ConnState,
	}

	serverStarted := make(chan struct{})
	serverErrors := make(chan error, 1)
	go func() {
		close(serverStarted)
		if err := mockServer.Serve(mal); err != nil {
			serverErrors <- err
		}
	}()

	<-serverStarted
	its.Empty(serverErrors)

	_, mockServerPort, err := net.SplitHostPort(mal.Addr().String())
	its.Nil(err)

	mockServerURL := fmt.Sprintf("http://127.0.0.1:%s", mockServerPort)

	dialer := new(net.Dialer)
	var dialCalls int32
	transport := new(http.Transport)
	transport.DialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
		atomic.AddInt32(&dialCalls, 1)
		return dialer.DialContext(ctx, network, addr)
	}

	client := http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodGet, mockServerURL, nil)
	its.Nil(err)

	itsFine := func(res *http.Response, err error) {
		its.Nil(err)
		its.Equal(http.StatusOK, res.StatusCode)
		_, readErr := io.Copy(ioutil.Discard, res.Body)
		its.Nil(readErr)
		its.Nil(res.Body.Close())
	}

	itsFine(client.Do(req))
	its.Equal(1, dialCalls)
	its.Equal(1, serverCalls)

	itsFine(client.Do(req))
	its.Equal(1, dialCalls)
	its.Equal(2, serverCalls)

	<-time.After(defaultMaxConnAge + time.Millisecond)

	itsFine(client.Do(req))
	its.Equal(1, dialCalls)
	its.Equal(3, serverCalls)

	itsFine(client.Do(req))
	its.Equal(2, dialCalls)
	its.Equal(4, serverCalls)
}
