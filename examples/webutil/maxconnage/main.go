/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"expvar"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

var (
	requests = new(expvar.Int)
)

func bindAddr() string {
	if value := os.Getenv("BIND_ADDR"); value != "" {
		return value
	}
	return "127.0.0.1:0"
}

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

var (
	statusOK = []byte(`{"status":"ok"}`)
)

func handler(rw http.ResponseWriter, req *http.Request) {
	defer requests.Add(1)
	rw.Header().Set(contentType, applicationJSON)
	rw.WriteHeader(http.StatusOK)
	rw.Write(statusOK)
}

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	log := logger.Prod()
	defer log.Close()

	ln, err := net.Listen("tcp", bindAddr())
	if err != nil {
		logger.MaybeFatal(log, err)
		exitCode = 1
		return
	}

	maxAgeListener := webutil.NewMaxAgeListener(ln, 5*time.Second, 10*time.Second)

	go async.NewInterval(func(_ context.Context) error {
		log.Infof("--- stats ---")
		log.Infof("requests: %d", requests.Value())
		log.Infof("opened: %d", maxAgeListener.ConnsOpened.Value())
		log.Infof("header closed: %d", maxAgeListener.ConnsHeaderClosed.Value())
		log.Infof("force closed: %d", maxAgeListener.ConnsForcedClosed.Value())
		return nil
	}, 10*time.Second).Start()

	log.Infof("listening on: %s", ln.Addr().String())

	server := http.Server{
		ErrorLog: logger.StdlibShim(log),
	}
	maxAgeListener.ApplyServer(&server)

	if err := server.Serve(maxAgeListener); err != nil {
		logger.MaybeFatal(log, err)
		exitCode = 1
		return
	}
}
