package email

import (
	"context"
)

// Sender is a generalized sender.
type Sender interface {
	Send(context.Context, Message) error
}
