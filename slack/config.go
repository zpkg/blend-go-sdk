package slack

import "github.com/blend/go-sdk/configutil"

// Config represents the required fields for the config.
type Config struct {
	Username  string `json:"username,omitempty" yaml:"username,omitempty" env:"SLACK_USERNAME"`
	Channel   string `json:"channel,omitempty" yaml:"channel,omitempty" env:"SLACK_CHANNEL"`
	IconURL   string `json:"iconURL,omitempty" yaml:"iconURL,omitempty" env:"SLACK_ICON_URL"`
	IconEmoji string `json:"iconEmoji,omitempty" yaml:"iconEmoji,omitempty" env:"SLACK_ICON_EMOJI"`
	Webhook   string `json:"webhook,omitempty" yaml:"webhook,omitempty"  env:"SLACK_WEBHOOK"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.Channel) == 0 && len(c.Webhook) == 0
}

// UsernameOrDefault returns a property or a default.
func (c Config) UsernameOrDefault(inherited ...string) string {
	return configutil.CoalesceString(c.Username, "", inherited...)
}

// ChannelOrDefault returns a property or default.
func (c Config) ChannelOrDefault(inherited ...string) string {
	return configutil.CoalesceString(c.Channel, "", inherited...)
}

// IconURLOrDefault returns a property or default.
func (c Config) IconURLOrDefault(inherited ...string) string {
	return configutil.CoalesceString(c.IconURL, "", inherited...)
}

// IconEmojiOrDefault returns a property or default.
func (c Config) IconEmojiOrDefault(inherited ...string) string {
	return configutil.CoalesceString(c.IconEmoji, "", inherited...)
}

// WebhookOrDefault returns the webhook url.
func (c Config) WebhookOrDefault(defaults ...string) string {
	return configutil.CoalesceString(c.Webhook, "", defaults...)
}
