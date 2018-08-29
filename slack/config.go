package slack

import "github.com/blend/go-sdk/util"

// Config represents the required fields for the config.
type Config struct {
	WebhookURL string `json:"webhookURL,omitempty" yaml:"webhookURL,omitempty"`
}

// GetWebhookURL returns the webhook url.
func (c Config) GetWebhookURL(defaults ...string) string {
	return util.Coalesce.String(c.WebhookURL, "", defaults...)
}
