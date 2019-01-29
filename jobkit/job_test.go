package jobkit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/uuid"
)

func TestJobProperties(t *testing.T) {
	assert := assert.New(t)

	job := NewJob(func(ctx context.Context) error {
		return nil
	})
	assert.NotNil(job.action)

	assert.NotEmpty(job.Name())
	job.WithName("foo")
	assert.Equal("foo", job.Name())

	assert.Nil(job.Schedule())
	job.WithSchedule(cron.EverySecond())
	assert.NotNil(job.Schedule())

	assert.Nil(job.NotificationsConfig())
	job.WithNotificationsConfig(&NotificationsConfig{})
	assert.NotNil(job.NotificationsConfig())

	assert.Zero(job.Timeout())
	job.WithTimeout(time.Second)
	assert.Equal(time.Second, job.Timeout())
}

func TestJobLifecycleHooksNotificationsUnset(t *testing.T) {
	assert := assert.New(t)

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:   uuid.V4().String(),
		Name: "test-job",
	})

	slackMessages := make(chan slack.Message, 1)

	job := &Job{
		slackClient: slack.MockWebhookSender(slackMessages),
	}

	job.OnStart(ctx)
	assert.Empty(slackMessages)

	job.OnComplete(ctx)
	assert.Empty(slackMessages)

	job.OnFailure(ctx)
	assert.Empty(slackMessages)

	job.OnCancellation(ctx)
	assert.Empty(slackMessages)

	job.OnBroken(ctx)
	assert.Empty(slackMessages)

	job.OnFixed(ctx)
	assert.Empty(slackMessages)
}

func TestJobLifecycleHooksNotificationsSetDisabled(t *testing.T) {
	assert := assert.New(t)

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:   uuid.V4().String(),
		Name: "test-job",
	})

	slackMessages := make(chan slack.Message, 1)

	job := &Job{
		slackClient: slack.MockWebhookSender(slackMessages),
		notifications: &NotificationsConfig{
			NotifyOnStart:   OptBool(false),
			NotifyOnSuccess: OptBool(false),
			NotifyOnFailure: OptBool(false),
			NotifyOnBroken:  OptBool(false),
			NotifyOnFixed:   OptBool(false),
		},
	}

	job.OnStart(ctx)
	assert.Empty(slackMessages)

	job.OnComplete(ctx)
	assert.Empty(slackMessages)

	job.OnFailure(ctx)
	assert.Empty(slackMessages)

	job.OnCancellation(ctx)
	assert.Empty(slackMessages)

	job.OnBroken(ctx)
	assert.Empty(slackMessages)

	job.OnFixed(ctx)
	assert.Empty(slackMessages)
}

func TestJobLifecycleHooksNotificationsSetEnabled(t *testing.T) {
	assert := assert.New(t)

	ctx := cron.WithJobInvocation(context.Background(), &cron.JobInvocation{
		ID:   uuid.V4().String(),
		Name: "test-job",
		Err:  fmt.Errorf("only a test"),
	})

	slackMessages := make(chan slack.Message, 6)

	job := &Job{
		slackClient: slack.MockWebhookSender(slackMessages),
		notifications: &NotificationsConfig{
			NotifyOnStart:   OptBool(true),
			NotifyOnSuccess: OptBool(true),
			NotifyOnFailure: OptBool(true),
			NotifyOnBroken:  OptBool(true),
			NotifyOnFixed:   OptBool(true),
		},
	}

	job.OnStart(ctx)
	job.OnComplete(ctx)
	job.OnFailure(ctx)
	job.OnCancellation(ctx)
	job.OnBroken(ctx)
	job.OnFixed(ctx)

	assert.Len(slackMessages, 6)

	msg := <-slackMessages
	assert.Contains("cron.started", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.complete", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.failed", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.cancelled", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.broken", msg.Text)

	msg = <-slackMessages
	assert.Contains("cron.fixed", msg.Text)
}
