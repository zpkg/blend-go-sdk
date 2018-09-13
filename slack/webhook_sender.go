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
		Config:        cfg,
	}
}

// WebhookSender sends slack webhooks.
type WebhookSender struct {
	*webutil.RequestSender
	Config *Config
}

// ApplyDefaults applies defaults.
func (whs WebhookSender) ApplyDefaults(message Message) Message {
	if len(message.Username) == 0 && whs.Config != nil {
		message.Username = whs.Config.GetUsername()
	}
	if len(message.Channel) == 0 && whs.Config != nil {
		message.Channel = whs.Config.GetChannel()
	}
	if len(message.IconURL) == 0 && whs.Config != nil {
		message.IconURL = whs.Config.GetIconURL()
	}
	if len(message.IconEmoji) == 0 && whs.Config != nil {
		message.IconEmoji = whs.Config.GetIconEmoji()
	}

	return message
}

// Send sends a slack hook.
func (whs WebhookSender) Send(message Message) error {
	res, err := whs.SendJSON(whs.ApplyDefaults(message))
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
