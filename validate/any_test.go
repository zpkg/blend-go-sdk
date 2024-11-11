/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/ex"
)

func TestForbidden(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any(nil).Forbidden()()
	assert.Nil(verr)

	verr = Any((*string)(nil)).Forbidden()()
	assert.Nil(verr)

	verr = Any("foo").Forbidden()()
	assert.NotNil(verr)
	assert.Equal(ErrForbidden, ErrCause(verr))
}

func TestRequired(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").Required()()
	assert.Nil(verr)

	verr = Any(nil).Required()()
	assert.NotNil(verr)
	assert.Equal(ErrRequired, ErrCause(verr))

	verr = Any((*string)(nil)).Required()()
	assert.NotNil(verr)
	assert.Equal(ErrRequired, ErrCause(verr))
}

func TestNotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").NotZero()()
	assert.Nil(verr)

	verr = Any(nil).NotZero()()
	assert.NotNil(verr)
	assert.Equal(ErrNotZero, ErrCause(verr))

	verr = Any((*string)(nil)).NotZero()()
	assert.NotNil(verr)
	assert.Equal(ErrNotZero, ErrCause(verr))
}

func TestAnyRefZero(t *testing.T) {
	assert := assert.New(t)

	type zeroTest struct {
		ID    int
		Value string
	}

	testCases := [...]struct {
		Input    interface{}
		Expected error
	}{
		{
			Input:    nil,
			Expected: nil,
		},
		{
			Input:    (*string)(nil),
			Expected: nil,
		},
		{
			Input:    0,
			Expected: nil,
		},
		{
			Input:    1,
			Expected: ErrZero,
		},
		{
			Input:    "",
			Expected: nil,
		},
		{
			Input:    "foo",
			Expected: ErrZero,
		},
		{
			Input:    zeroTest{},
			Expected: nil,
		},
		{
			Input:    zeroTest{ID: 2},
			Expected: ErrZero,
		},
	}

	for index, tc := range testCases {
		verr := Any(tc.Input).Zero()()
		assert.Equal(tc.Expected, ErrCause(verr), index)
	}
}

func TestAnyRefEmpty(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    interface{}
		Expected error
	}{
		{
			Input:    nil,
			Expected: ErrNonLengthType,
		},
		{
			Input:    0,
			Expected: ErrNonLengthType,
		},
		{
			Input:    []string{},
			Expected: nil,
		},
		{
			Input:    ([]string)(nil),
			Expected: nil,
		},
		{
			Input:    map[string]interface{}{},
			Expected: nil,
		},
		{
			Input:    (map[string]interface{})(nil),
			Expected: nil,
		},
		{
			Input:    "",
			Expected: nil,
		},
		{
			Input:    make(chan struct{}),
			Expected: nil,
		},
		{
			Input:    (chan struct{})(nil),
			Expected: nil,
		},
		{
			Input:    []string{"a", "b"},
			Expected: ErrEmpty,
		},
		{
			Input:    map[string]int{"hi": 1},
			Expected: ErrEmpty,
		},
		{
			Input:    "foo",
			Expected: ErrEmpty,
		},
	}

	for index, tc := range testCases {
		verr := Any(tc.Input).Empty()()
		assert.Equal(tc.Expected, ErrCause(verr), index)
	}
}

func TestAnyRefNotEmpty(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    interface{}
		Expected error
	}{
		{
			Input:    nil,
			Expected: ErrNonLengthType,
		},
		{
			Input:    0,
			Expected: ErrNonLengthType,
		},
		{
			Input:    []string{},
			Expected: ErrNotEmpty,
		},
		{
			Input:    ([]string)(nil),
			Expected: ErrNotEmpty,
		},
		{
			Input:    map[string]interface{}{},
			Expected: ErrNotEmpty,
		},
		{
			Input:    (map[string]interface{})(nil),
			Expected: ErrNotEmpty,
		},
		{
			Input:    "",
			Expected: ErrNotEmpty,
		},
		{
			Input:    make(chan struct{}),
			Expected: ErrNotEmpty,
		},
		{
			Input:    (chan struct{})(nil),
			Expected: ErrNotEmpty,
		},
		{
			Input:    []string{"a", "b"},
			Expected: nil,
		},
		{
			Input:    map[string]int{"hi": 1},
			Expected: nil,
		},
		{
			Input:    "foo",
			Expected: nil,
		},
	}

	for index, tc := range testCases {
		verr := Any(tc.Input).NotEmpty()()
		assert.Equal(tc.Expected, ErrCause(verr), index)
	}
}

func TestAnyRefLen(t *testing.T) {
	assert := assert.New(t)

	err := Any(1234).Len(10)()
	assert.NotNil(err)
	assert.Equal(ErrNonLengthType, ex.ErrClass(err))

	verr := Any([]int{1, 2, 3, 4}).Len(4)()
	assert.Nil(verr)

	verr = Any(map[int]bool{1: true, 2: true}).Len(2)()
	assert.Nil(verr)

	verr = Any([]int{}).Len(4)()
	assert.NotNil(verr)
	assert.Equal(ErrLen, ErrCause(verr))
}

func TestAnyRefNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any(nil).Nil()()
	assert.Nil(verr)

	var nilPtr *string
	verr = Any(nilPtr).Nil()()
	assert.Nil(verr)

	verr = Any("foo").Nil()()
	assert.NotNil(verr)
	assert.Equal(ErrNil, ErrCause(verr))
}

func TestAnyRefNotNil(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").NotNil()()
	assert.Nil(verr)

	verr = Any(nil).NotNil()()
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ErrCause(verr))

	var nilPtr *string
	verr = Any(nilPtr).NotNil()()
	assert.NotNil(verr)
	assert.Equal(ErrNotNil, ErrCause(verr))
}

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").Equals("foo")()
	assert.Nil(verr)

	verr = Any(nil).Equals(nil)()
	assert.Nil(verr)

	verr = Any("foo").Equals("bar")()
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ErrCause(verr))

	verr = Any(nil).Equals("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrEquals, ErrCause(verr))
}

func TestAnyRefNotEquals(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").NotEquals("bar")()
	assert.Nil(verr)

	verr = Any(nil).NotEquals("foo")()
	assert.Nil(verr)

	verr = Any("foo").NotEquals("foo")()
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ErrCause(verr))

	verr = Any(nil).NotEquals(nil)()
	assert.NotNil(verr)
	assert.Equal(ErrNotEquals, ErrCause(verr))
}

func TestAnyRefAllow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("foo").Allow("foo", "bar", "baz")()
	assert.Nil(verr)
	verr = Any("bar").Allow("foo", "bar", "baz")()
	assert.Nil(verr)
	verr = Any("baz").Allow("foo", "bar", "baz")()
	assert.Nil(verr)

	verr = Any("what").Allow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrAllowed, ErrCause(verr))
}

func TestAnyRefDisallow(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Any("what").Disallow("foo", "bar", "baz")()
	assert.Nil(verr)

	verr = Any("foo").Disallow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ErrCause(verr))
	verr = Any("bar").Disallow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ErrCause(verr))
	verr = Any("baz").Disallow("foo", "bar", "baz")()
	assert.NotNil(verr)
	assert.Equal(ErrDisallowed, ErrCause(verr))
}
