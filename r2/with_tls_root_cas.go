package r2

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

// WithTLSRootCAs sets the client tls root ca pool.
func WithTLSRootCAs(pool *x509.CertPool) Option {
	return func(r *Request) {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			if typed.TLSClientConfig == nil {
				typed.TLSClientConfig = &tls.Config{}
			}
			typed.TLSClientConfig.RootCAs = pool
		}
	}
}
