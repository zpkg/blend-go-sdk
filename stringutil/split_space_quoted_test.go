package stringutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitSpaceQuoted(t *testing.T) {
	assert := assert.New(t)

	values := SplitSpaceQuoted("")
	assert.Len(values, 0)

	values = SplitSpaceQuoted("foo")
	assert.Len(values, 1)
	assert.Equal("foo", values[0])

	values = SplitSpaceQuoted("foo bar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo  bar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo\tbar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo \tbar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo bar  ")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo bar baz")
	assert.Len(values, 3)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])
	assert.Equal("baz", values[2])

	values = SplitSpaceQuoted(`foo "bar baz"`)
	assert.Len(values, 2, fmt.Sprintf("%#v", values))
	assert.Equal("foo", values[0])
	assert.Equal(`"bar baz"`, values[1])

	values = SplitSpaceQuoted(`foo --config="bar baz"`)
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal(`--config="bar baz"`, values[1])

	values = SplitSpaceQuoted(`foo --config='bar baz="hi"'`)
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal(`--config='bar baz="hi"'`, values[1])

	values = SplitSpaceQuoted(`foo --config="bar baz='hi'"`)
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal(`--config="bar baz='hi'"`, values[1])
}
