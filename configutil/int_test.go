/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestInt(t *testing.T) {
	assert := assert.New(t)

	intValue := Int(0)
	ptr, err := intValue.Int(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	intValue = Int(1234)
	ptr, err = intValue.Int(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(1234, *ptr)
}
