package webutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMustParseURL(t *testing.T) {
	assert := assert.New(t)

	output := MustParseURL("https://foo.bar.com/bailey?buzz=muzz")
	assert.NotNil(output)
	assert.Equal("https", output.Scheme)
	assert.Equal("foo.bar.com", output.Host)
	assert.Equal("/bailey", output.Path)
	assert.NotEmpty(output.Query())

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		output = MustParseURL(":not-a-url-at-all")
	}()
	assert.NotNil(err)
}

func TestURLWithScheme(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/bailey?buzz=muzz")
	assert.Equal("http", URLWithScheme(original, "http").Scheme)
}

func TestURLWithHost(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/bailey?buzz=muzz")
	assert.Equal("blend.com", URLWithHost(original, "blend.com").Host)
}

func TestURLWithPath(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/bailey?buzz=muzz")
	assert.Equal("not-bailey", URLWithPath(original, "not-bailey").Path)
}

func TestURLWithRawQuery(t *testing.T) {
	assert := assert.New(t)

	original := MustParseURL("https://foo.bar.com/bailey?buzz=muzz")
	assert.Equal("dog=cool", URLWithRawQuery(original, "dog=cool").RawQuery)
}
