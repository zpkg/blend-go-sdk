package slack

import (
	"io/ioutil"

	"net/http"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/webutil"
)

const (
	// ErrNon200 is the exception class when a non-200 is returned from slack.
	ErrNon200 = "slack; non-200 status code returned from remote"
)

// NewWebhookSender creates a new webhook sender.
func NewWebhookSender(cfg *Config) *WebhookSender {
	return &WebhookSender{
		RequestSender: webutil.NewRequestSender(webutil.MustParseURL(cfg.GetWebhook())),
	}
}

// WebhookSender sends slack webhooks.
type WebhookSender struct {
	*webutil.RequestSender
}

// Send sends a slack hook.
func (whs WebhookSender) Send(message *Message) error {
	res, err := whs.SendJSON(message)
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
