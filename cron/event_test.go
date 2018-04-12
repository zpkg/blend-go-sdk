package cron

import (
	"bytes"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
	logger "github.com/blend/go-sdk/logger"
)

func TestEventStartedListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	textBuffer := bytes.NewBuffer(nil)
	jsonBuffer := bytes.NewBuffer(nil)
	all := logger.New().WithFlags(logger.AllFlags()).
		WithRecoverPanics(false).
		WithWriter(logger.NewTextWriter(textBuffer)).
		WithWriter(logger.NewJSONWriter(jsonBuffer))

	defer all.Close()

	all.Listen(FlagStarted, "default", NewEventListener(func(e *Event) {
		defer wg.Done()

		assert.Equal(FlagStarted, e.Flag())
		assert.False(e.Timestamp().IsZero())
		assert.Equal("test_task", e.TaskName())
		assert.False(e.Complete())
		assert.Nil(e.Err())
		assert.Zero(e.Elapsed())
	}))

	go func() { all.Trigger(&Event{flag: FlagStarted, taskName: "test_task"}) }()
	go func() { all.Trigger(&Event{flag: FlagStarted, taskName: "test_task"}) }()
	wg.Wait()
	all.Drain()

	assert.NotEmpty(textBuffer.String())
	assert.NotEmpty(jsonBuffer.String())
}

func TestEventInterfaces(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(FlagComplete, "test_task").WithHeading("heading").WithLabel("foo", "bar")

	eventProvider, isEvent := logger.MarshalEvent(e)
	assert.True(isEvent)
	assert.Equal(FlagComplete, eventProvider.Flag())
	assert.False(eventProvider.Timestamp().IsZero())

	headingProvider, isHeadingProvider := logger.MarshalEventHeading(e)
	assert.True(isHeadingProvider)
	assert.Equal("heading", headingProvider.Heading())

	metaProvider, isMetaProvider := logger.MarshalEventMeta(e)
	assert.True(isMetaProvider)
	assert.Equal("bar", metaProvider.Labels()["foo"])
}
