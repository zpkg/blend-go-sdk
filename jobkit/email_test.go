package jobkit

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/email"
)

func TestNewEmailMessage(t *testing.T) {
	assert := assert.New(t)

	message, err := NewEmailMessage(string(cron.FlagComplete), &cron.JobInvocation{JobName: "test", Elapsed: time.Millisecond}, email.WithFrom("jobkit@blend.com"), email.WithTo("foo@bar.com"), email.WithCC("baileydog@blend.com"))
	assert.Nil(err)
	assert.Equal("test :: cron.complete", message.Subject)
	assert.NotEmpty(message.From)
	assert.Equal("jobkit@blend.com", message.From)
	assert.NotEmpty(message.To)
	assert.Equal("foo@bar.com", message.To[0])
	assert.NotEmpty(message.CC)
	assert.Equal("baileydog@blend.com", message.CC[0])
	assert.NotEmpty(message.HTMLBody)
	assert.NotEmpty(message.TextBody)
}
