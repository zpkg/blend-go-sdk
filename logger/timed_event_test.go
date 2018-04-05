package logger

import (
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimedEventListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	all := New().WithFlags(AllFlags())
	defer all.Close()
	all.Listen(Flag("test-flag"), "default", NewTimedEventListener(func(te *TimedEvent) {
		defer wg.Done()
		assert.Equal("test-flag", te.Flag())
		assert.NotZero(te.Elapsed())
		assert.Equal("foo bar", te.Message())
	}))

	go func() { all.Trigger(Timedf(Flag("test-flag"), time.Millisecond, "foo %s", "bar")) }()
	go func() { all.Trigger(Timedf(Flag("test-flag"), time.Millisecond, "foo %s", "bar")) }()
	wg.Wait()
}

func TestTimedEventInterfaces(t *testing.T) {
	assert := assert.New(t)

	ee := Timedf(Fatal, time.Millisecond, "foo %s", "bar").WithHeading("heading").WithLabel("foo", "bar")

	eventProvider, isEvent := marshalEvent(ee)
	assert.True(isEvent)
	assert.Equal(Fatal, eventProvider.Flag())
	assert.False(eventProvider.Timestamp().IsZero())

	headingProvider, isHeadingProvider := marshalEventHeading(ee)
	assert.True(isHeadingProvider)
	assert.Equal("heading", headingProvider.Heading())

	metaProvider, isMetaProvider := marshalEventMeta(ee)
	assert.True(isMetaProvider)
	assert.Equal("bar", metaProvider.Labels()["foo"])
}
