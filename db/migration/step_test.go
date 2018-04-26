package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStep(t *testing.T) {
	assert := assert.New(t)

	step := NewStep(AlwaysRun(), NoOp)
	assert.True(step.TransactionBound())
	assert.NotNil(step.guard)
	assert.NotNil(step.body)
	assert.Empty(step.Label())
	assert.Nil(step.Parent())
	assert.Nil(step.Collector())

	step.WithLabel("test")
	assert.Equal("test", step.Label())

	step.WithParent(NewGroup())
	assert.NotNil(step.Parent())

	step.WithCollector(&Collector{})
	assert.NotNil(step.Collector())

	step.WithTransactionBound(false)
	assert.False(step.TransactionBound())
}
