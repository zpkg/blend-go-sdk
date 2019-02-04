package slack

// ApplyMessageOptions applies a set of options against a message and returns the mutated copy.
func ApplyMessageOptions(m Message, options ...MessageOption) Message {
	for _, option := range options {
		option(&m)
	}
	return m
}

// MessageOption is a mutator for messages.
type MessageOption func(m *Message)

// WithChannel sets the channel.
func WithChannel(channel string) MessageOption {
	return func(m *Message) {
		m.Channel = channel
	}
}

// WithChannelOrDefault sets the channel if its unset.
func WithChannelOrDefault(channel string) MessageOption {
	return func(m *Message) {
		if len(m.Channel) == 0 {
			m.Channel = channel
		}
	}
}

// WithResponseType sets the response type.
func WithResponseType(responseType string) MessageOption {
	return func(m *Message) {
		m.ResponseType = responseType
	}
}

// WithUsername sets the username.
func WithUsername(username string) MessageOption {
	return func(m *Message) {
		m.Username = username
	}
}

// WithUsernameOrDefault sets the username.
func WithUsernameOrDefault(username string) MessageOption {
	return func(m *Message) {
		if len(m.Username) == 0 {
			m.Username = username
		}
	}
}

// WithIconEmoji sets the icon emoji.
func WithIconEmoji(emoji string) MessageOption {
	return func(m *Message) {
		m.IconEmoji = emoji
	}
}

// WithIconEmojiOrDefault sets the icon emoji.
func WithIconEmojiOrDefault(emoji string) MessageOption {
	return func(m *Message) {
		if len(m.IconEmoji) == 0 {
			m.IconEmoji = emoji
		}
	}
}

// WithIconURL sets the icon url.
func WithIconURL(url string) MessageOption {
	return func(m *Message) {
		m.IconURL = url
	}
}

// WithIconURLOrDefault sets the icon url.
func WithIconURLOrDefault(url string) MessageOption {
	return func(m *Message) {
		if len(m.IconURL) == 0 {
			m.IconURL = url
		}
	}
}

// WithMessageAttachment adds a message attachment.
func WithMessageAttachment(attachment MessageAttachment) MessageOption {
	return func(m *Message) {
		m.Attachments = append(m.Attachments, attachment)
	}
}
