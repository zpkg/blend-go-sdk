package slack

import (
	"context"
	"io/ioutil"

	"net/http"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/webutil"
)

const (
	// ErrNon200 is the exception class when a non-200 is returned from slack.
	ErrNon200 = "slack; non-200 status code returned from remote"
)

var (
	_ Sender = (*WebhookSender)(nil)
)

// New creates a new webhook sender.
func New(cfg *Config) *WebhookSender {
	return &WebhookSender{
		RequestSender: webutil.NewRequestSender(webutil.MustParseURL(cfg.WebhookOrDefault())),
		Config:        cfg,
	}
}

// WebhookSender sends slack webhooks.
type WebhookSender struct {
	*webutil.RequestSender
	Config *Config
}

// Defaults returns default message options.
func (whs WebhookSender) Defaults() []MessageOption {
	return []MessageOption{
		WithUsernameOrDefault(whs.Config.UsernameOrDefault()),
		WithChannelOrDefault(whs.Config.ChannelOrDefault()),
		WithIconEmojiOrDefault(whs.Config.IconEmojiOrDefault()),
		WithIconURLOrDefault(whs.Config.IconURLOrDefault()),
	}
}

// Send sends a slack hook.
func (whs WebhookSender) Send(ctx context.Context, message Message) error {
	res, err := whs.SendJSON(ctx, ApplyMessageOptions(message, whs.Defaults()...))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return exception.New(err)
		}
		return exception.New(ErrNon200).WithMessage(string(contents))
	}
	return nil
}

// PostMessage posts a basic message to a given chanel.
func (whs WebhookSender) PostMessage(channel, messageText string, options ...MessageOption) error {
	message := Message{
		Channel: channel,
		Text:    messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.Send(context.Background(), message)
}

// PostMessageContext posts a basic message to a given chanel with a given context.
func (whs WebhookSender) PostMessageContext(ctx context.Context, channel, messageText string, options ...MessageOption) error {
	message := Message{
		Channel: channel,
		Text:    messageText,
	}
	for _, option := range options {
		option(&message)
	}
	return whs.Send(ctx, message)
}
