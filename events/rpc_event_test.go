package logger

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestRPCEventListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)

	textBuffer := bytes.NewBuffer(nil)
	jsonBuffer := bytes.NewBuffer(nil)
	all := New().WithFlags(AllFlags()).WithRecoverPanics(false).
		WithWriter(NewTextWriter(textBuffer)).
		WithWriter(NewJSONWriter(jsonBuffer))
	defer all.Close()

	all.Listen(RPC, "default", NewRPCEventListener(func(re *RPCEvent) {
		defer wg.Done()
		assert.Equal(RPC, re.Flag())
		assert.Equal("/test", re.Method())
		assert.NotZero(re.Elapsed())
	}))

	go func() {
		defer wg.Done()
		all.Trigger(NewRPCEvent("/test", time.Millisecond))
	}()
	go func() {
		defer wg.Done()
		all.Trigger(NewRPCEvent("/test", time.Millisecond))
	}()
	wg.Wait()
	all.Drain()

	assert.NotEmpty(textBuffer.String())
	assert.NotEmpty(jsonBuffer.String())
}

func TestRPCEventProperties(t *testing.T) {
	assert := assert.New(t)

	e := NewRPCEvent("", 0)

	assert.False(e.Timestamp().IsZero())
	assert.True(e.WithTimestamp(time.Time{}).Timestamp().IsZero())

	assert.Empty(e.Labels())
	assert.Equal("bar", e.WithLabel("foo", "bar").Labels()["foo"])

	assert.Empty(e.Annotations())
	assert.Equal("zar", e.WithAnnotation("moo", "zar").Annotations()["moo"])

	assert.Equal(RPC, e.Flag())
	assert.Equal(Error, e.WithFlag(Error).Flag())

	assert.Empty(e.Headings())
	assert.Equal([]string{"Heading"}, e.WithHeadings("Heading").Headings())

	assert.Empty(e.Engine())
	assert.Equal("grpc", e.WithEngine("grpc").Engine())

	assert.Empty(e.Method())
	assert.Equal("/test", e.WithMethod("/test").Method())

	assert.Zero(e.Elapsed())
	assert.Equal(time.Millisecond, e.WithElapsed(time.Millisecond).Elapsed())

	assert.Empty(e.ContentType())
	assert.Equal("content-type", e.WithContentType("content-type").ContentType())

	assert.Empty(e.Authority())
	assert.Equal("authority", e.WithAuthority("authority").Authority())

}
