/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestCreateServer(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	caKeyPair := KeyPair{
		Cert: string(caCertLiteral),
		Key:  string(caKeyLiteral),
	}
	authority, err := NewCertBundle(caKeyPair)
	assert.Nil(err)

	assert.Len(authority.CertificateDERs, 1)
	assert.Len(authority.Certificates, 1)

	server, err := CreateServer("warden-server", authority, OptDNSNames("warden-server-test"))
	assert.Nil(err)
	assert.Len(server.Certificates, 2)
	assert.Len(server.CertificateDERs, 2)
	assert.Len(server.Certificates[0].DNSNames, 1)
}
