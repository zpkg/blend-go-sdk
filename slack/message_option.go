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

// WithUsername sets the channel.
func WithUsername(username string) MessageOption {
	return func(m *Message) {
		m.Username = username
	}
}

// WithIconEmoji sets the icon emoji.
func WithIconEmoji(emoji string) MessageOption {
	return func(m *Message) {
		m.IconEmoji = emoji
	}
}

// WithIconURL sets the icon url.
func WithIconURL(url string) MessageOption {
	return func(m *Message) {
		m.IconURL = url
	}
}
