package email

// Message is a message to send via. ses.
type Message struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	CC       []string `json:"cc"`
	BCC      []string `json:"bcc"`
	Subject  string   `json:"subject"`
	TextBody string   `json:"textBody"`
	HTMLBody string   `json:"htmlBody"`
}

// IsZero returns if the object is set or not.
func (m Message) IsZero() bool {
	return len(m.To) == 0
}
