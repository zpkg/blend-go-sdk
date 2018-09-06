package slack

import "github.com/blend/go-sdk/util"

// Config represents the required fields for the config.
type Config struct {
	Username  string `json:"username,omitempty" yaml:"username,omitempty"`
	Channel   string `json:"channel,omitempty" yaml:"channel,omitempty"`
	IconURL   string `json:"iconURL,omitempty" yaml:"iconURL,omitempty"`
	IconEmoji string `json:"iconEmoji,omitempty" yaml:"iconEmoji,omitempty"`
	Webhook   string `json:"webhook,omitempty" yaml:"webhook,omitempty"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.Webhook) == 0
}

// GetUsername returns a property or a default.
func (c Config) GetUsername(inherited ...string) string {
	return util.Coalesce.String(c.Username, "", inherited...)
}

// GetChannel returns a property or default.
func (c Config) GetChannel(inherited ...string) string {
	return util.Coalesce.String(c.Channel, "", inherited...)
}

// GetIconURL returns a property or default.
func (c Config) GetIconURL(inherited ...string) string {
	return util.Coalesce.String(c.IconURL, "", inherited...)
}

// GetIconEmoji returns a property or default.
func (c Config) GetIconEmoji(inherited ...string) string {
	return util.Coalesce.String(c.IconEmoji, "", inherited...)
}

// GetWebhook returns the webhook url.
func (c Config) GetWebhook(defaults ...string) string {
	return util.Coalesce.String(c.Webhook, "", defaults...)
}
