package jobkit

import (
	"context"
	"fmt"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/slack"
)

// SendSlackMessage sends a complete slack message.
func SendSlackMessage(ctx context.Context, sender slack.Sender, jobStatus JobStatus, ji *cron.JobInvocation, options ...slack.MessageOption) error {
	if sender == nil {
		return nil
	}
	return sender.Send(ctx, NewSlackMessage(ji.Name, jobStatus, ji.Err, ji.Elapsed, options...))
}

// NewSlackMessage returns a new job started message.
func NewSlackMessage(jobName string, jobStatus JobStatus, err error, elapsed time.Duration, options ...slack.MessageOption) slack.Message {
	message := slack.Message{}
	if err != nil {
		message.Attachments = append(message.Attachments,
			slack.MessageAttachment{
				Text:  fmt.Sprintf("%s %s", jobName, jobStatus),
				Color: "#ff0000",
			},
			slack.MessageAttachment{
				Text:  fmt.Sprintf("error: %+v", err),
				Color: "#ff0000",
			},
		)
	} else {
		message.Attachments = append(message.Attachments,
			slack.MessageAttachment{
				Text:  fmt.Sprintf("%s %s", jobName, jobStatus),
				Color: "#00ff00",
			},
		)
	}

	if elapsed > 0 {
		message.Attachments = append(message.Attachments,
			slack.MessageAttachment{
				Text: fmt.Sprintf("%v elapsed", elapsed),
			},
		)
	}

	return slack.ApplyMessageOptions(message, options...)
}
