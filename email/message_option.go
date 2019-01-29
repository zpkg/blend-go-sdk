package email

// ApplyMessageOptions applies options to a message and returns
// the mutated copy.
func ApplyMessageOptions(m Message, options ...MessageOption) Message {
	for _, option := range options {
		option(&m)
	}
	return m
}

// MessageOption is a mutator for messages.
type MessageOption func(m *Message)

// WithFrom sets the from address for a message.
func WithFrom(from string) MessageOption {
	return func(m *Message) {
		m.From = from
	}
}

// WithTo sets the to address for a message.
func WithTo(to ...string) MessageOption {
	return func(m *Message) {
		m.To = to
	}
}

// WithCC sets the cc addresses for a message.
func WithCC(cc ...string) MessageOption {
	return func(m *Message) {
		m.CC = cc
	}
}

// WithBCC sets the bcc addresses for a message.
func WithBCC(bcc ...string) MessageOption {
	return func(m *Message) {
		m.BCC = bcc
	}
}

// WithSubject sets the subject for a message.
func WithSubject(subject string) MessageOption {
	return func(m *Message) {
		m.Subject = subject
	}
}
