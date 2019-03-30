package airbrake

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestFrames(t *testing.T) {
	assert := assert.New(t)

	ex := exception.As(exception.New("this is a test"))
	fr := frames(ex.Stack)
	assert.NotEmpty(fr, fmt.Sprintf("%T", ex.Stack))

	fr = frames(exception.As(exception.New("this is a test", exception.OptStack(exception.StackStrings([]string{"foo", "bar"})))).Stack)
	assert.Empty(fr)
}
