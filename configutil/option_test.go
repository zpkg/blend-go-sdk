package configutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func Test_OptAddPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPaths("foo", "bar")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"foo", "bar"}, options.Paths)

	assert.Nil(OptAddPaths("moo", "loo")(&options))
	assert.Len(options.Paths, 4)
	assert.Equal([]string{"foo", "bar", "moo", "loo"}, options.Paths)
}

func Test_OptAddPreferredPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPreferredPaths("foo", "bar")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"foo", "bar"}, options.Paths)

	assert.Nil(OptAddPreferredPaths("moo", "loo")(&options))
	assert.Len(options.Paths, 4)
	assert.Equal([]string{"moo", "loo", "foo", "bar"}, options.Paths)
}

func Test_OptPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptPaths("foo", "bar")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"foo", "bar"}, options.Paths)

	assert.Nil(OptPaths("moo", "loo")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"moo", "loo"}, options.Paths)
}

func Test_OptEnv(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Empty(options.Env)
	assert.Nil(OptEnv(env.Vars{"FOO": "bar"})(&options))
	assert.NotEmpty(options.Env)
	assert.Len(options.Env, 1)
	assert.Equal("bar", options.Env["FOO"])
}
