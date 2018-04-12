package logger

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestErrorEventListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	textBuffer := bytes.NewBuffer(nil)
	jsonBuffer := bytes.NewBuffer(nil)
	all := New().WithFlags(AllFlags()).WithRecoverPanics(false).
		WithWriter(NewTextWriter(textBuffer)).
		WithWriter(NewJSONWriter(jsonBuffer))
	defer all.Close()

	all.Listen(Fatal, "default", NewErrorEventListener(func(e *ErrorEvent) {
		defer wg.Done()

		assert.Equal(Fatal, e.Flag())
		assert.Equal("foo bar", e.Err().Error())
	}))

	go func() { all.Trigger(NewErrorEvent(Fatal, fmt.Errorf("foo bar"))) }()
	go func() { all.Trigger(NewErrorEvent(Fatal, fmt.Errorf("foo bar"))) }()
	wg.Wait()
	all.Drain()

	assert.NotEmpty(textBuffer.String())
	assert.NotEmpty(jsonBuffer.String())
}

func TestErrorEventInterfaces(t *testing.T) {
	assert := assert.New(t)

	ee := NewErrorEvent(Fatal,
		exception.New("this is a test").
			WithMessagef("this is a message").
			WithStack(exception.StackStrings([]string{"foo", "bar"}))).
		WithHeading("heading").
		WithLabel("foo", "bar")

	eventProvider, isEvent := MarshalEvent(ee)
	assert.True(isEvent)
	assert.Equal(Fatal, eventProvider.Flag())
	assert.False(eventProvider.Timestamp().IsZero())

	headingProvider, isHeadingProvider := MarshalEventHeading(ee)
	assert.True(isHeadingProvider)
	assert.Equal("heading", headingProvider.Heading())

	metaProvider, isMetaProvider := MarshalEventMeta(ee)
	assert.True(isMetaProvider)
	assert.Equal("bar", metaProvider.Labels()["foo"])
}

func TestErrorEventProperties(t *testing.T) {
	assert := assert.New(t)

	ee := NewErrorEvent(Fatal, nil)
	assert.False(ee.Timestamp().IsZero())
	assert.True(ee.WithTimestamp(time.Time{}).Timestamp().IsZero())

	assert.Empty(ee.Labels())
	assert.Equal("bar", ee.WithLabel("foo", "bar").Labels()["foo"])

	assert.Empty(ee.Annotations())
	assert.Equal("zar", ee.WithAnnotation("moo", "zar").Annotations()["moo"])

	assert.Equal(Fatal, ee.Flag())
	assert.Equal(Error, ee.WithFlag(Error).Flag())

	assert.Empty(ee.Heading())
	assert.Equal("Heading", ee.WithHeading("Heading").Heading())

	assert.Nil(ee.Err())
	assert.Equal(fmt.Errorf("foo"), ee.WithErr(fmt.Errorf("foo")).Err())

	assert.Nil(ee.State())
	assert.Equal("State", ee.WithState("State").State())

	assert.Empty(ee.FlagTextColor())
	assert.Equal(ColorWhite, ee.WithFlagTextColor(ColorWhite).FlagTextColor())
}
