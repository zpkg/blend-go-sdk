package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStep(t *testing.T) {
	assert := assert.New(t)

	step := NewStep(AlwaysRun(), NoOp)
	assert.NotNil(step.guard)
	assert.NotNil(step.body)
	assert.Empty(step.Label())

	step.WithLabel("test")
	assert.Equal("test", step.Label())
}
