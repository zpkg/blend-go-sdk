/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestStdlibShim(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	log, err := New(
		OptOutput(buf),
		OptAll(),
		OptText(OptTextHideTimestamp(), OptTextNoColor()),
	)
	assert.Nil(err)
	defer log.Close()

	shim := StdlibShim(log, OptShimWriterEventProvider(ShimWriterErrorEventProvider("error")))

	shim.Println("this is a test")
	shim.Println("this is another test")

	assert.NotEmpty(buf.String())
	assert.Equal("[error] this is a test\n[error] this is another test\n", buf.String())
}
