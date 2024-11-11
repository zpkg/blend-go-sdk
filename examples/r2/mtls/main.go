/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"os"

	"github.com/zpkg/blend-go-sdk/certutil"
	"github.com/zpkg/blend-go-sdk/graceful"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/r2"
	"github.com/zpkg/blend-go-sdk/web"
	"github.com/zpkg/blend-go-sdk/webutil"
)

func fatal(log logger.FatalReceiver, err error) {
	log.Fatal(err)
	os.Exit(1)
}

func main() {
	log := logger.All()

	// create the ca
	ca, err := certutil.CreateCertificateAuthority()
	if err != nil {
		fatal(log, err)
	}

	caKeyPair, err := ca.GenerateKeyPair()
	if err != nil {
		fatal(log, err)
	}

	caPool, err := ca.CertPool()
	if err != nil {
		fatal(log, err)
	}

	// create the server certs
	server, err := certutil.CreateServer("mtls-example.local", ca, certutil.OptSubjectAlternateNames("localhost"))
	if err != nil {
		fatal(log, err)
	}
	serverKeyPair, err := server.GenerateKeyPair()
	if err != nil {
		fatal(log, err)
	}

	client, err := certutil.CreateClient("mtls-client", ca)
	if err != nil {
		fatal(log, err)
	}
	clientKeyPair, err := client.GenerateKeyPair()
	if err != nil {
		fatal(log, err)
	}

	serverCertManager, err := certutil.NewCertManagerWithKeyPairs(serverKeyPair, []certutil.KeyPair{caKeyPair}, clientKeyPair)
	if err != nil {
		fatal(log, err)
	}

	// create a server
	app, err := web.New(
		web.OptLog(log),
		web.OptBindAddr("127.0.0.1:5000"),
		web.OptTLSConfig(serverCertManager.TLSConfig),
		web.OptServerOptions(
			webutil.OptHTTPServerErrorLog(
				logger.StdlibShim(log, logger.OptShimWriterEventProvider(
					logger.ShimWriterMessageEventProvider("http.error"),
				)),
			),
		),
	)
	if err != nil {
		fatal(log, err)
	}

	go func() {
		<-app.NotifyStarted()

		// make some requests ...
		log.Info("making a secure request")

		if _, err := r2.New("https://localhost:5000",
			r2.OptTLSRootCAs(caPool),
			r2.OptTLSClientCert([]byte(clientKeyPair.Cert), []byte(clientKeyPair.Key))).Discard(); err != nil {
			fatal(log, err)
		} else {
			log.Info("secure request success")
		}

		log.Info("making an insecure request")
		if _, err := r2.New("https://localhost:5000", r2.OptTLSRootCAs(caPool)).Discard(); err != nil {
			log.Error(err)
		}
	}()

	if err := graceful.Shutdown(app); err != nil {
		fatal(log, err)
	}
}
