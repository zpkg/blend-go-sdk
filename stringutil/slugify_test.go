package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSlugify(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", Slugify("foo"))
	assert.Equal("foo-bar", Slugify("foo bar"))
	assert.Equal("foo-bar", Slugify("foo\tbar"))
	assert.Equal("foo-bar", Slugify("foo\nbar"))
	assert.Equal("foo-bar-ba-", Slugify("foo bar ba/"))
}
