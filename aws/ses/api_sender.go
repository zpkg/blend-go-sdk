package ses

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	awsSes "github.com/aws/aws-sdk-go/service/ses"
	"github.com/blend/go-sdk/exception"
)

// APISender is an aws ses email sender.
type APISender struct {
	session *session.Session
	client  *awsSes.SES
}

// Send sends a message.
func (s *APISender) Send(ctx context.Context, m Message) error {
	if s.client == nil {
		return nil
	}
	input := &awsSes.SendEmailInput{
		Source: &m.From,
		Destination: &awsSes.Destination{
			ToAddresses:  strPtrs(m.To),
			CcAddresses:  strPtrs(m.CC),
			BccAddresses: strPtrs(m.BCC),
		},
		Message: &awsSes.Message{
			Subject: &awsSes.Content{
				Charset: &defaultCharset,
				Data:    &m.Subject,
			},
			Body: &awsSes.Body{},
		},
	}

	if len(m.HTMLBody) > 0 {
		input.Message.Body.Html = &awsSes.Content{
			Charset: &defaultCharset,
			Data:    &m.HTMLBody,
		}
	}

	if len(m.TextBody) > 0 {
		input.Message.Body.Text = &awsSes.Content{
			Charset: &defaultCharset,
			Data:    &m.TextBody,
		}
	}

	_, err := s.client.SendEmailWithContext(ctx, input)
	return exception.New(err)
}

func strPtrs(values []string) []*string {
	output := make([]*string, len(values))
	for i := 0; i < len(values); i++ {
		output[i] = &values[i]
	}
	return output
}
