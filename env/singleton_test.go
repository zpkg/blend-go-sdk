/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env_test

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/env"
)

func TestClear(t *testing.T) {
	assert := assert.New(t)

	vars := env.Vars{
		"Foo": "bar",
	}
	env.SetEnv(vars)
	assert.NotEmpty(env.Env())

	env.Clear()
	assert.Empty(env.Env())
}
