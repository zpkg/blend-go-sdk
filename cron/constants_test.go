package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestErrors(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsJobNotLoaded(ErrJobNotLoaded))
	assert.True(IsJobNotLoaded(exception.NewFromErr(ErrJobNotLoaded)))
	assert.False(IsJobNotLoaded(ErrJobAlreadyLoaded))
	assert.False(IsJobNotLoaded(exception.NewFromErr(ErrJobAlreadyLoaded)))

	assert.True(IsJobAlreadyLoaded(ErrJobAlreadyLoaded))
	assert.True(IsJobAlreadyLoaded(exception.NewFromErr(ErrJobAlreadyLoaded)))
	assert.False(IsJobAlreadyLoaded(ErrJobNotLoaded))
	assert.False(IsJobAlreadyLoaded(exception.NewFromErr(ErrJobNotLoaded)))

	assert.True(IsTaskNotFound(ErrTaskNotFound))
	assert.True(IsTaskNotFound(exception.NewFromErr(ErrTaskNotFound)))
	assert.False(IsTaskNotFound(ErrJobNotLoaded))
	assert.False(IsTaskNotFound(exception.NewFromErr(ErrJobNotLoaded)))
}
