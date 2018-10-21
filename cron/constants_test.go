package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestErrors(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsJobNotLoaded(ErrJobNotLoaded))
	assert.True(IsJobNotLoaded(exception.New(ErrJobNotLoaded)))
	assert.False(IsJobNotLoaded(ErrJobAlreadyLoaded))
	assert.False(IsJobNotLoaded(exception.New(ErrJobAlreadyLoaded)))

	assert.True(IsJobAlreadyLoaded(ErrJobAlreadyLoaded))
	assert.True(IsJobAlreadyLoaded(exception.New(ErrJobAlreadyLoaded)))
	assert.False(IsJobAlreadyLoaded(ErrJobNotLoaded))
	assert.False(IsJobAlreadyLoaded(exception.New(ErrJobNotLoaded)))

	assert.True(IsJobNotFound(ErrJobNotFound))
	assert.True(IsJobNotFound(exception.New(ErrJobNotFound)))
	assert.False(IsJobNotFound(ErrJobNotLoaded))
	assert.False(IsJobNotFound(exception.New(ErrJobNotLoaded)))
}
