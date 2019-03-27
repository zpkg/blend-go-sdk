package certutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"time"

	"github.com/blend/go-sdk/exception"
)

// CreateCA creates a ca cert bundle.
// The cert bundle can be used to generate client and server certificates.
// It creates a private key for the csr.
func CreateCA(options ...CertOption) (*CertBundle, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, exception.New(err)
	}
	return CreateCAFromPrivateKey(pk, options...)
}

// CreateCAFromPrivateKey creates a ca cert bundle from a given private key.
// The cert bundle can be used to generate client and server certificates.
func CreateCAFromPrivateKey(pk *rsa.PrivateKey, options ...CertOption) (*CertBundle, error) {
	var output CertBundle
	var err error
	output.PrivateKey = pk
	output.PublicKey = &pk.PublicKey
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, exception.New(err)
	}

	csr := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().UTC().AddDate(DefaultCANotAfterYears, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	for _, option := range options {
		option(&csr)
	}

	der, err := x509.CreateCertificate(rand.Reader, &csr, &csr, output.PublicKey, output.PrivateKey)
	if err != nil {
		return nil, exception.New(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, exception.New(err)
	}
	output.CertificateDERs = [][]byte{der}
	output.Certificates = []x509.Certificate{*cert}
	return &output, nil
}
