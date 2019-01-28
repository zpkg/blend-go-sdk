package slack

import "context"

// NewMockWebhookSender creates a new mock sender.
func NewMockWebhookSender() MockWebhookSender {
	return MockWebhookSender(make(chan Message))
}

// MockWebhookSender is a mocked sender.
type MockWebhookSender chan Message

// Send sends a mocked message.
func (ms MockWebhookSender) Send(ctx context.Context, m Message) error {
	ms <- m
	return nil
}
