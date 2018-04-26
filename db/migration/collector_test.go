package migration

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestCollectorUnset(t *testing.T) {
	assert := assert.New(t)

	// none of these should panic.
	var unset *Collector
	unset.Applyf(nil, "")
	unset.Skipf(nil, "")
	unset.Error(nil, nil)

	assert.True(true, "should have reached here")
}

func TestCollectorApplyf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	log := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	c := &Collector{output: log}

	c.Applyf(NewStep(AlwaysRun(), NoOp), "a good one")

	assert.Equal("[db.migration] -- applied -- a good one\n", output.String())
	assert.Equal(1, c.applied)
	assert.Equal(0, c.skipped)
	assert.Equal(0, c.failed)
	assert.Equal(1, c.total)
}

func TestCollectorSkipf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	log := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	c := &Collector{output: log}

	c.Skipf(NewStep(AlwaysRun(), NoOp), "a good one")

	assert.Equal("[db.migration] -- skipped -- a good one\n", output.String())
	assert.Equal(0, c.applied)
	assert.Equal(1, c.skipped)
	assert.Equal(0, c.failed)
	assert.Equal(1, c.total)
}

func TestCollectorError(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	log := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	c := &Collector{output: log}

	assert.NotNil(c.Error(NewStep(AlwaysRun(), NoOp), fmt.Errorf("a good one")))

	assert.Equal("[db.migration] -- failed -- a good one\n", output.String())
	assert.Equal(0, c.applied)
	assert.Equal(0, c.skipped)
	assert.Equal(1, c.failed)
	assert.Equal(1, c.total)
}

func TestCollectorLabels(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	log := logger.New().WithFlags(logger.AllFlags()).WithWriter(logger.NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	c := &Collector{output: log}

	c.Applyf(NewStep(AlwaysRun(), NoOp).WithLabel("test label"), "a good one")

	assert.Equal("[db.migration] -- applied test label -- a good one\n", output.String())
}
