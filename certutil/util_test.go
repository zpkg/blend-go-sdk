package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestReadFiles(t *testing.T) {
	assert := assert.New(t)

	files, err := ReadFiles("testdata/client.cert.pem", "testdata/client.key.pem")
	assert.Nil(err)
	assert.Len(files, 2)
}

func TestParseCertPEM(t *testing.T) {
	assert := assert.New(t)

	certs, err := ParseCertPEM(certLiteral)
	assert.Nil(err)
	assert.Len(certs, 2)
	assert.Equal(certLiteralCommonName, certs[0].Subject.CommonName)
}

func TestCommonNamesForCertPair(t *testing.T) {
	assert := assert.New(t)

	kp := KeyPair{Cert: string(certLiteral)}

	certPEM, err := kp.CertBytes()
	assert.Nil(err)

	commonNames, err := CommonNamesForCertPEM(certPEM)
	assert.Nil(err)
	assert.NotEmpty(commonNames)
	assert.Equal(certLiteralCommonName, commonNames[0])
}
