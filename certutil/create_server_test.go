package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	ca, err := CreateCA()
	assert.Nil(err)

	server, err := CreateServer("warden-serevr", &ca)
	assert.Nil(err)
	assert.Len(server.Certificates, 2)
	assert.Len(server.CertificateDERs, 2)
}
