package slack

// Response is a slack response
type Response struct {
	OK              bool            `json:"ok"`
	Channel         string          `json:"channel,omitempty"`
	ThreadTimestamp string          `json:"ts,omitempty"`
	Message         ResponseMessage `json:"message,omitempty"`
	Error           string          `json:"error,omitempty"`
}

// ResponseMessage is a message received from a slack response
type ResponseMessage struct {
	Text            string              `json:"text,omitempty"`
	Username        string              `json:"username,omitempty"`
	BotID           string              `json:"bot_id,omitempty"`
	Attachments     []MessageAttachment `json:"attachments"`
	Type            string              `json:"type,omitempty"`
	SubType         string              `json:"bot_message,omitempty"`
	ThreadTimestamp string              `json:"ts,omitempty"`
}
