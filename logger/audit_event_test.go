package logger

import (
	"sync"
	"testing"
	"time"

	assert "github.com/blend/go-sdk/assert"
)

func TestAuditEventListener(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	all := New().WithFlags(AllFlags()).WithRecoverPanics(false)
	defer all.Close()

	all.Listen(Audit, "default", NewAuditEventListener(func(e *AuditEvent) {
		defer wg.Done()

		assert.Equal(Audit, e.Flag())
		assert.Equal("principal", e.Principal())
		assert.Equal("verb", e.Verb())
	}))

	go func() { all.Trigger(NewAuditEvent("principal", "verb", "noun")) }()
	go func() { all.Trigger(NewAuditEvent("principal", "verb", "noun")) }()
	wg.Wait()

}

func TestAuditEventInterfaces(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent("principal", "verb", "noun").WithHeading("heading").WithLabel("foo", "bar")

	eventProvider, isEvent := marshalEvent(ae)
	assert.True(isEvent)
	assert.Equal(Audit, eventProvider.Flag())
	assert.False(eventProvider.Timestamp().IsZero())

	headingProvider, isHeadingProvider := marshalEventHeading(ae)
	assert.True(isHeadingProvider)
	assert.Equal("heading", headingProvider.Heading())

	metaProvider, isMetaProvider := marshalEventMeta(ae)
	assert.True(isMetaProvider)
	assert.Equal("bar", metaProvider.Labels()["foo"])
}

func TestAuditEventProperties(t *testing.T) {
	assert := assert.New(t)

	ae := NewAuditEvent("principal", "verb", "noun")
	assert.False(ae.Timestamp().IsZero())
	assert.Equal("principal", ae.Principal())
	assert.Equal("verb", ae.Verb())
	assert.Equal("noun", ae.Noun())

	assert.True(ae.WithTimestamp(time.Time{}).Timestamp().IsZero())

	assert.Empty(ae.Heading())
	assert.Equal("Heading", ae.WithHeading("Heading").Heading())

	assert.Empty(ae.Subject())
	assert.Equal("Subject", ae.WithSubject("Subject").Subject())

	assert.Empty(ae.Property())
	assert.Equal("Property", ae.WithProperty("Property").Property())

	assert.Empty(ae.RemoteAddress())
	assert.Equal("RemoteAddress", ae.WithRemoteAddress("RemoteAddress").Property())
}
