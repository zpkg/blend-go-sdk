package main

import (
	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All()

	// create the ca
	ca, err := certutil.CreateCA()
	if err != nil {
		log.SyncFatalExit(err)
	}

	caKeyPair, err := ca.KeyPair()
	if err != nil {
		log.SyncFatalExit(err)
	}

	caPool, err := ca.CertPool()
	if err != nil {
		log.SyncFatalExit(err)
	}

	// create the server certs
	server, err := certutil.CreateServer("mtls-example.local", &ca, "localhost")
	if err != nil {
		log.SyncFatalExit(err)
	}
	serverKeyPair, err := server.KeyPair()
	if err != nil {
		log.SyncFatalExit(err)
	}

	client, err := certutil.CreateClient("mtls-client", &ca)
	if err != nil {
		log.SyncFatalExit(err)
	}
	clientKeyPair, err := client.KeyPair()
	if err != nil {
		log.SyncFatalExit(err)
	}

	serverCertManager, err := certutil.NewCertManagerWithKeyPairs(serverKeyPair, []certutil.KeyPair{caKeyPair}, clientKeyPair)
	if err != nil {
		log.SyncFatalExit(err)
	}

	// create a server
	app := web.New().WithLogger(log).WithBindAddr("127.0.0.1:5000")
	app.WithTLSConfig(serverCertManager.TLSConfig)
	go func() {
		if err := graceful.Shutdown(app); err != nil {
			log.SyncFatalExit(err)
		}
	}()
	<-app.NotifyStarted()

	// make some requests ...

	log.SyncInfof("making a secure request")
	if err := r2.New("https://localhost:5000",
		r2.TLSRootCAs(caPool),
		r2.TLSClientCert([]byte(clientKeyPair.Cert), []byte(clientKeyPair.Key)),
	).Discard(); err != nil {
		log.SyncFatalExit(err)
	} else {
		log.SyncInfof("secure request success")
	}

	log.SyncInfof("making an insecure request")
	if err := r2.New("https://localhost:5000", r2.TLSRootCAs(caPool)).Discard(); err != nil {
		log.SyncFatalExit(err)
	}
}
