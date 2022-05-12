/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCertFileWatcher(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	tempKey, err := ioutil.TempFile("", "")
	assert.Nil(err)
	defer func() {
		os.Remove(tempKey.Name())
	}()

	tempCert, err := ioutil.TempFile("", "")
	assert.Nil(err)
	defer func() {
		os.Remove(tempCert.Name())
	}()

	_, err = tempKey.Write(keyLiteral)
	assert.Nil(err)

	_, err = tempCert.Write(certLiteral)
	assert.Nil(err)

	assert.Nil(tempKey.Close())
	assert.Nil(tempCert.Close())

	w, err := NewCertFileWatcher(tempCert.Name(), tempKey.Name())
	assert.Nil(err)
	assert.NotNil(w.Certificate)

	assert.Nil(w.Reload())
	assert.NotNil(w.Certificate)
}
