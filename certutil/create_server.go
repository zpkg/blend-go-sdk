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

// CreateServer creates a ca cert bundle.
func CreateServer(commonName string, ca *CertBundle, options ...CertOption) (*CertBundle, error) {
	if ca == nil {
		return nil, exception.New("must provide a ca cert bundle")
	}

	var output CertBundle
	var err error
	output.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, exception.New(err)
	}
	output.PublicKey = &output.PrivateKey.PublicKey

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, exception.New(err)
	}
	csr := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Warden"},
			Country:      []string{"United States"},
			Province:     []string{"California"},
			Locality:     []string{"San Francisco"},
		},
		NotBefore:   time.Now().UTC(),
		NotAfter:    time.Now().UTC().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	csr.DNSNames = []string{commonName}
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
