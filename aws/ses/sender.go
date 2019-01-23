package ses

import (
	"context"

	awsSes "github.com/aws/aws-sdk-go/service/ses"
	"github.com/blend/go-sdk/aws"
)

var (
	defaultCharset = "UTF-8"
)

// New returns a new sender.
func New(cfg *aws.Config) Sender {
	return &APISender{
		client: awsSes.New(aws.NewSession(cfg)),
	}
}

// Sender is a generalized sender.
type Sender interface {
	Send(context.Context, Message) error
}
