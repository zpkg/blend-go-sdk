/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package datadog

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_RateSampler(t *testing.T) {
	assert := assert.New(t)

	sampler := RateSampler(0.25)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.True(passed > 225, fmt.Sprint(passed))
	assert.True(passed < 280, fmt.Sprint(passed))
}

func Test_RateSampler_FullOn(t *testing.T) {
	assert := assert.New(t)

	sampler := RateSampler(1)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Equal(passed, 1024)
}

func Test_RateSampler_FullOff(t *testing.T) {
	assert := assert.New(t)

	sampler := RateSampler(0)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Zero(passed)
}

func Test_RateSampler_FullOff_Negative(t *testing.T) {
	assert := assert.New(t)

	sampler := RateSampler(-1)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Zero(passed)
}
