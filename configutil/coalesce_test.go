package configutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCoalesceString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", CoalesceString("", ""))
	assert.Equal("foo", CoalesceString("", "foo"))
	assert.Equal("bar", CoalesceString("", "foo", "bar"))
	assert.Equal("bar", CoalesceString("", "foo", "bar", "baz"))
	assert.Equal("moo", CoalesceString("moo", "foo", "bar", "baz"))
}

func TestCoalesceBool(t *testing.T) {
	assert := assert.New(t)

	assert.False(CoalesceBool(nil, false))
	assert.True(CoalesceBool(nil, true))
	assert.False(CoalesceBool(nil, true, false))
	assert.True(CoalesceBool(refBool(true), false, false))
}

func refBool(value bool) *bool {
	return &value
}

func TestCoalesceInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CoalesceInt(0, 1))
	assert.Equal(2, CoalesceInt(0, 1, 2))
	assert.Equal(2, CoalesceInt(0, 1, 2, 3))
	assert.Equal(4, CoalesceInt(4, 1, 2, 3))
}

func TestCoalesceInt32(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CoalesceInt32(0, 1))
	assert.Equal(2, CoalesceInt32(0, 1, 2))
	assert.Equal(2, CoalesceInt32(0, 1, 2, 3))
	assert.Equal(4, CoalesceInt32(4, 1, 2, 3))
}

func TestCoalesceInt64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CoalesceInt64(0, 1))
	assert.Equal(2, CoalesceInt64(0, 1, 2))
	assert.Equal(2, CoalesceInt64(0, 1, 2, 3))
	assert.Equal(4, CoalesceInt64(4, 1, 2, 3))
}

func TestCoalesceFloat32(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CoalesceFloat32(0, 1))
	assert.Equal(2, CoalesceFloat32(0, 1, 2))
	assert.Equal(2, CoalesceFloat32(0, 1, 2, 3))
	assert.Equal(4, CoalesceFloat32(4, 1, 2, 3))
}

func TestCoalesceFloat64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CoalesceFloat64(0, 1))
	assert.Equal(2, CoalesceFloat64(0, 1, 2))
	assert.Equal(2, CoalesceFloat64(0, 1, 2, 3))
	assert.Equal(4, CoalesceFloat64(4, 1, 2, 3))
}

func TestCoalesceDuration(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CoalesceDuration(0, 1))
	assert.Equal(2, CoalesceDuration(0, 1, 2))
	assert.Equal(2, CoalesceDuration(0, 1, 2, 3))
	assert.Equal(4, CoalesceDuration(4, 1, 2, 3))
}

func TestCoalesceTime(t *testing.T) {
	assert := assert.New(t)

	zero := time.Time{}
	one := time.Date(2018, 01, 01, 12, 00, 00, 00, time.UTC)
	two := time.Date(2018, 02, 02, 12, 00, 00, 00, time.UTC)
	three := time.Date(2018, 03, 03, 12, 00, 00, 00, time.UTC)
	four := time.Date(2018, 04, 04, 12, 00, 00, 00, time.UTC)

	assert.Equal(one, CoalesceTime(zero, one))
	assert.Equal(two, CoalesceTime(zero, one, two))
	assert.Equal(two, CoalesceTime(zero, one, two, three))
	assert.Equal(four, CoalesceTime(four, one, two, three))
}

func TestCoalesceStrings(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CoalesceStrings(nil, nil))
	assert.NotEmpty(CoalesceStrings(nil, []string{"foo"}))
	assert.NotEmpty(CoalesceStrings(nil, []string{}, []string{"bar"}))
	assert.NotEmpty(CoalesceStrings([]string{"moo"}, []string{}, []string{}))
}

func TestCoalesceBytes(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CoalesceBytes(nil, nil))
	assert.NotEmpty(CoalesceBytes(nil, []byte("foo")))
	assert.NotEmpty(CoalesceBytes(nil, []byte{}, []byte("bar")))
	assert.NotEmpty(CoalesceBytes([]byte("moo"), []byte{}, []byte{}))
}
