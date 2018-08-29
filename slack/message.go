package slack

// Message is a message sent to slack.
type Message struct {
	ResponseType string `json:"response_type,omitempty"`
	Text         string `json:"text"`
	Username     string `json:"username,omitempty"`
	UnfurlLinks  bool   `json:"unfurl_links"`
	IconEmoji    string `json:"icon_emoji,omitempty"`

	Attachments []MessageAttachment `json:"attachments"`
}

// MessageAttachment is an attachment for a message.
type MessageAttachment struct {
	Title      string                   `json:"title,omitempty"`
	Color      string                   `json:"color,omitempty"`
	Pretext    string                   `json:"pretext,omitempty"`
	Text       string                   `json:"text,omitempty"`
	MarkdownIn []string                 `json:"mrkdwn_in,omitempty"`
	Fields     []MessageAttachmentField `json:"fields,omitempty"`
}

// MessageAttachmentField is a field on an attachment.
type MessageAttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short"`
}
