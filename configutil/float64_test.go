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

func TestFloat64(t *testing.T) {
	assert := assert.New(t)

	floatValue := Float64(0)
	ptr, err := floatValue.Float64(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	floatValue = Float64(3.14)
	ptr, err = floatValue.Float64(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(3.14, *ptr)
}
