package certutil

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/blend/go-sdk/exception"
)

// NewClientConfig returns a new client config.
func NewClientConfig(clientCert KeyPair, certificateAuthorities []KeyPair) (*tls.Config, error) {
	clientCertPEM, err := clientCert.CertBytes()
	if err != nil {
		return nil, exception.New(err)
	}
	clientKeyPEM, err := clientCert.KeyBytes()
	if err != nil {
		return nil, exception.New(err)
	}

	if len(clientCertPEM) == 0 || len(clientKeyPEM) == 0 {
		return nil, exception.New("empty cert or key pem")
	}
	cert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, exception.New(err)
	}
	config := &tls.Config{}

	rootCAPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, exception.New(err)
	}

	for _, caCert := range certificateAuthorities {
		contents, err := caCert.CertBytes()
		if err != nil {
			return nil, exception.New(err)
		}
		if ok := rootCAPool.AppendCertsFromPEM(contents); !ok {
			return nil, exception.New("failed to append ca cert file")
		}
	}

	config.Certificates = []tls.Certificate{cert}
	config.RootCAs = rootCAPool
	return config, nil
}
