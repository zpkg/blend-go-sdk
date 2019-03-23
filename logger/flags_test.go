package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestEventFlagSetEnable(t *testing.T) {
	assert := assert.New(t)

	set := NewFlags().WithEnabled("FOO")
	set.Enable("TEST")
	assert.True(set.IsEnabled("TEST"))
	assert.True(set.IsEnabled("FOO"))
	assert.False(set.IsEnabled("NOT_TEST"))
}

func TestEventFlagSetDisable(t *testing.T) {
	assert := assert.New(t)

	set := NewFlags()
	set.Enable("TEST")
	assert.True(set.IsEnabled("TEST"))
	set.Disable("TEST")
	assert.False(set.IsEnabled("TEST"))
}

func TestEventFlagSetEnableAll(t *testing.T) {
	assert := assert.New(t)

	set := NewFlags()
	set.SetAll()
	assert.True(set.IsEnabled("TEST"))
	assert.True(set.IsEnabled("NOT_TEST"))
	assert.True(set.IsEnabled("NOT_TEST"))
	set.Disable("TEST")
	assert.True(set.IsEnabled("NOT_TEST"))
	assert.False(set.IsEnabled("TEST"))
}

func TestEventFlagSetFromEnvironment(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarEventFlags, "error,info,web.request")
	defer env.Env().Restore(EnvVarEventFlags)

	set := NewFlagsFromEnv()
	assert.True(set.IsEnabled(Error))
	assert.True(set.IsEnabled(Info))
	assert.False(set.IsEnabled(Fatal))
}

func TestEventFlagSetFromEnvironmentWithDisabled(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarEventFlags, "all,-debug")
	defer env.Env().Restore(EnvVarEventFlags)

	set := NewFlagsFromEnv()
	assert.True(set.IsEnabled(Error))
	assert.True(set.IsEnabled(Fatal))
	assert.True(set.IsEnabled(Flag("foo")))
	assert.False(set.IsEnabled(Debug))
}

func TestEventFlagSetFromEnvironmentAll(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarEventFlags, "all")
	defer env.Env().Restore(EnvVarEventFlags)

	set := NewFlagsFromEnv()
	assert.True(set.All())
	assert.False(set.None())
	assert.True(set.IsEnabled(Error))
}

func TestEventFlagSetFromEnvNone(t *testing.T) {
	assert := assert.New(t)

	env.Env().Set(EnvVarEventFlags, "none")
	defer env.Env().Restore(EnvVarEventFlags)

	set := NewFlagsFromEnv()
	assert.False(set.All())
	assert.True(set.None())
	assert.False(set.IsEnabled(Error))
}

func TestEventFlagNoneEnableEvents(t *testing.T) {
	assert := assert.New(t)

	flags := NewFlagsNone()
	assert.False(flags.IsEnabled("test_flag"))

	flags.Enable("test_flag")
	assert.True(flags.IsEnabled("test_flag"))
}

func TestEventSetCoalesceWith(t *testing.T) {
	assert := assert.New(t)

	first := NewFlags(Info)
	first.CoalesceWith(NewFlags(Warning))
	assert.True(first.IsEnabled(Info))
	assert.True(first.IsEnabled(Warning))
	assert.False(first.IsEnabled(Fatal))

	second := NewFlags(Info)
	second.CoalesceWith(NewFlagsFromValues("-info"))
	assert.False(second.IsEnabled(Info))
}

func TestFlagSetNone(t *testing.T) {
	assert := assert.New(t)

	assert.True(NoneFlags().None())
}

func TestNewHiddenFlagSetFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()
	env.Env().Set(EnvVarHiddenEventFlags, "debug,silly")

	set := NewHiddenFlagSetFromEnv()
	assert.False(set.IsEnabled(Info))
	assert.True(set.IsEnabled(Debug))
	assert.True(set.IsEnabled(Silly))
}

func TestFlagSetString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("none", NoneFlags().String())
	assert.Equal("all", AllFlags().String())

	fs := NewFlags(Info, Debug, Error)
	assert.Contains(fs.String(), "info")
	assert.Contains(fs.String(), "debug")
	assert.Contains(fs.String(), "error")

	nfs := AllFlags().WithDisabled(Fatal)
	assert.Equal("all, -fatal", nfs.String())
}
