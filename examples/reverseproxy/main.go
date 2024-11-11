/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"net"
	"net/http"
	"os"

	"github.com/zpkg/blend-go-sdk/graceful"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/r2"
	"github.com/zpkg/blend-go-sdk/reverseproxy"
	"github.com/zpkg/blend-go-sdk/webutil"
)

func main() {
	log := logger.Prod()

	_, err := r2.New("https://www.google.com").CopyTo(os.Stdout)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	upstream := reverseproxy.NewUpstream(webutil.MustParseURL("https://www.google.com"))
	proxy, _ := reverseproxy.NewProxy(
		reverseproxy.OptProxyUpstream(upstream),
		reverseproxy.OptProxySetHeaderValue(webutil.HeaderXForwardedProto, webutil.SchemeHTTPS),
		reverseproxy.OptProxySetHeaderValue(webutil.HeaderXForwardedHost, "www.google.com"),
		reverseproxy.OptProxyLog(log),
	)

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	server := http.Server{
		Handler: proxy,
	}
	log.Infof("listening on: %s", listener.Addr().String())
	if err := graceful.Shutdown(webutil.NewGracefulHTTPServer(&server, webutil.OptGracefulHTTPServerListener(listener))); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
