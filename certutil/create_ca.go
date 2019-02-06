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

// CreateCA creates a ca cert bundle.
// The cert bundle can be used to generate client and server certificates.
func CreateCA(options ...CertOption) (output CertBundle, err error) {
	output.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = exception.New(err)
		return
	}
	output.PublicKey = &output.PrivateKey.PublicKey

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		err = exception.New(err)
		return
	}

	csr := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "warden-ca",
			Organization: []string{"Warden"},
			Country:      []string{"United States"},
			Province:     []string{"California"},
			Locality:     []string{"San Francisco"},
		},
		NotBefore:             time.Now().UTC(),
		NotAfter:              time.Now().UTC().AddDate(10, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	for _, option := range options {
		option(&csr)
	}

	der, err := x509.CreateCertificate(rand.Reader, &csr, &csr, output.PublicKey, output.PrivateKey)
	if err != nil {
		err = exception.New(err)
		return
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		err = exception.New(err)
		return
	}
	output.CertificateDERs = [][]byte{der}
	output.Certificates = []x509.Certificate{*cert}
	return
}
