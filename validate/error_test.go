/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"fmt"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/ex"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	verr := Error(fmt.Errorf("this is a test"), nil)
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Empty(ex.ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), ErrCause(verr))

	verr = Error(fmt.Errorf("this is a test"), nil, "foo", "bar")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal("foobar", ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), ErrCause(verr))
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)

	verr := Errorf(fmt.Errorf("this is a test"), "foo", "minimum: %d", 30)
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal("minimum: 30", ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), ErrCause(verr))
}

func TestErrCause(t *testing.T) {
	assert := assert.New(t)

	err := ex.New(ErrNonLengthType)
	assert.Equal(ErrNonLengthType, ex.ErrClass(err))
	assert.Equal(ErrNonLengthType, ErrCause(err))

	verr := Error(ErrEmpty, "foo")
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrEmpty, ErrCause(verr))
}
