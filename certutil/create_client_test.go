/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/uuid"
)

func TestCreateClient(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	caKeyPair := KeyPair{
		Cert: string(caCertLiteral),
		Key:  string(caKeyLiteral),
	}
	ca, err := NewCertBundle(caKeyPair)
	assert.Nil(err)

	uid := uuid.V4().String()
	client, err := CreateClient(uid, ca)
	assert.Nil(err)
	assert.Len(client.Certificates, 2)
	assert.Len(client.CertificateDERs, 2)
}
