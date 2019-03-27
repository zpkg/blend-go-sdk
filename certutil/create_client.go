package certutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/blend/go-sdk/exception"
)

/*
CreateClient creates a client cert bundle associated with a given common name.

The CA must be passed in as a CertBundle.

Example:

	ca, err := certutil.NewCertBundle(certutil.KeyPairFromPaths("ca.crt", "ca.key"))
	if err != nil {
		return err
	}
	client, err := CreateClient("foo.bar.com", ca)
*/
func CreateClient(commonName string, ca *CertBundle, options ...CertOption) (*CertBundle, error) {
	if ca == nil {
		return nil, exception.New("must provide a ca cert bundle")
	}
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, exception.New(err)
	}
	return CreateClientFromPrivateKey(pk, commonName, ca, options...)
}

/*
CreateClientFromPrivateKey creates a client cert bundle associated with a given common name with a given private key.

The CA must be passed in as a CertBundle.

Example:

	pk, err := certutil.ReadPrivateKeyPEMFromPath("id_rsa")
	if err != nil {
		return err
	}
	ca, err := certutil.NewCertBundle(certutil.KeyPairFromPaths("ca.crt", "ca.key"))
	if err != nil {
		return err
	}
	client, err := CreateClientFromPrivateKey(pk, "foo.bar.com", ca)
*/
func CreateClientFromPrivateKey(pk *rsa.PrivateKey, commonName string, ca *CertBundle, options ...CertOption) (*CertBundle, error) {
	var output CertBundle
	var err error
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, exception.New(err)
	}
	csr := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:   time.Now().UTC(),
		NotAfter:    time.Now().UTC().AddDate(DefaultClientNotAfterYears, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	for _, option := range options {
		option(&csr)
	}

	der, err := x509.CreateCertificate(rand.Reader, &csr, &ca.Certificates[0], output.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, exception.New(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, exception.New(err)
	}
	output.CertificateDERs = append([][]byte{der}, ca.CertificateDERs...)
	output.Certificates = append([]x509.Certificate{*cert}, ca.Certificates...)
	return &output, nil
}
