package slack

import "github.com/blend/go-sdk/util"

// Config represents the required fields for the config.
type Config struct {
	WebhookURL string `json:"webhookURL,omitempty" yaml:"webhookURL,omitempty"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.WebhookURL) == 0
}

// GetWebhookURL returns the webhook url.
func (c Config) GetWebhookURL(defaults ...string) string {
	return util.Coalesce.String(c.WebhookURL, "", defaults...)
}
