package email

import (
	"context"
)

var (
	defaultCharset = "UTF-8"
)

// Sender is a generalized sender.
type Sender interface {
	Send(context.Context, Message) error
}
