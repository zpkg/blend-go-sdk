/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"crypto/rand"
	"crypto/x509"

	"github.com/zpkg/blend-go-sdk/ex"
)

// CreateCertificateAuthority creates a ca cert bundle from a given set of options.
// The cert bundle can be used to generate client and server certificates.
func CreateCertificateAuthority(options ...CertOption) (*CertBundle, error) {
	createOptions := DefaultOptionsCertificateAuthority

	if err := ResolveCertOptions(&createOptions, options...); err != nil {
		return nil, nil
	}

	var output CertBundle
	output.PrivateKey = createOptions.PrivateKey
	output.PublicKey = &createOptions.PrivateKey.PublicKey
	der, err := x509.CreateCertificate(rand.Reader, &createOptions.Certificate, &createOptions.Certificate, output.PublicKey, output.PrivateKey)
	if err != nil {
		return nil, ex.New(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, ex.New(err)
	}
	output.CertificateDERs = [][]byte{der}
	output.Certificates = []x509.Certificate{*cert}
	return &output, nil
}
