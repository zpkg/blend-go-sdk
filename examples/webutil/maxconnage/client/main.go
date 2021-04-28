/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"expvar"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
)

var (
	requests         = new(expvar.Int)
	requestsComplete = new(expvar.Int)
)

func main() {
	log := logger.Prod()

	transport := &http.Transport{
		DisableKeepAlives:   false,
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
	}

	ctx, cancel := context.WithCancel(graceful.Background())
	errors := make(chan error, runtime.NumCPU())

	start := time.Now()
	for x := 0; x < runtime.NumCPU(); x++ {
		fire(ctx, transport, errors)
	}

	select {
	case <-ctx.Done():
		log.Info("canceled!")
	case <-time.After(30 * time.Second):
		cancel()
		log.Info("done!")
	case err := <-errors:
		cancel()
		log.Info("errored!")
		log.Fatal(err)
	}

	fmt.Printf("%d/%d requests completed in %v\n", requestsComplete.Value(), requests.Value(), time.Since(start))
}

func fire(ctx context.Context, transport *http.Transport, errors chan error) {
	go func() {
		client := http.Client{
			Transport: transport,
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8081", nil)
		if err != nil {
			errors <- err
			return
		}
		var res *http.Response
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			requests.Add(1)
			res, err = client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "request error: %+v\n", err)
				errors <- err
				return
			}
			_, err = io.Copy(ioutil.Discard, res.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "request read error: %+v\n", err)
				errors <- err
				return
			}
			if err = res.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "request body close error: %+v\n", err)
				errors <- err
				return
			}
			requestsComplete.Add(1)
		}
	}()
}
