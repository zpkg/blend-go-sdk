package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCompressWhitespace(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", CompressWhitespace(""))
	assert.Equal("", CompressWhitespace(" "))
	assert.Equal("", CompressWhitespace("\n"))
	assert.Equal("", CompressWhitespace("\t"))

	assert.Equal("foo", CompressWhitespace(" foo"))
	assert.Equal("foo", CompressWhitespace("foo "))
	assert.Equal("foo", CompressWhitespace("foo\n"))

	assert.Equal("foo bar", CompressWhitespace("foo bar"))
	assert.Equal("foo bar", CompressWhitespace("foo\tbar"))
	assert.Equal("foo bar", CompressWhitespace("foo\nbar"))

	assert.Equal("foo bar", CompressWhitespace("foo  bar"))
	assert.Equal("foo bar", CompressWhitespace("foo\t\tbar"))
	assert.Equal("foo bar", CompressWhitespace("foo\n\nbar"))

	assert.Equal("foo bar baz", CompressWhitespace("foo  bar   baz"))
	assert.Equal("foo bar baz", CompressWhitespace("foo\t\t\tbar baz\n"))
	assert.Equal("foo bar baz", CompressWhitespace("foo\n\n\nbar\tbaz"))
}
