/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "bar",
	}
	assert.True(Equals{Key: "foo", Value: "far"}.Matches(valid))
	assert.False(Equals{Key: "zoo", Value: "buzz"}.Matches(valid))
	assert.False(Equals{Key: "foo", Value: "bar"}.Matches(valid))

	assert.Equal("foo == bar", Equals{Key: "foo", Value: "bar"}.String())
}
