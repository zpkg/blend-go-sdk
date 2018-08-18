package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewGroup(t *testing.T) {
	assert := assert.New(t)

	g := NewGroup(NewStep(AlwaysRun(), NoOp))
	assert.Empty(g.Label())
	assert.Len(g.steps, 1)

	g.WithLabel("test")
	assert.Equal("test", g.Label())
}
