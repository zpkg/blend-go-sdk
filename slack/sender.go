package slack

import "context"

// Sender is a type that can send slack messages.
type Sender interface {
	Send(context.Context, Message) error
}
