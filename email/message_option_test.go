package email

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMessageOption(t *testing.T) {
	assert := assert.New(t)

	m := Message{}

	assert.Empty(m.From)
	WithFrom("foo@bar.com")(&m)
	assert.Equal("foo@bar.com", m.From)

	assert.Empty(m.To)
	WithTo("buzz@bar.com", "fuzz@bar.com")(&m)
	assert.Equal([]string{"buzz@bar.com", "fuzz@bar.com"}, m.To)

	assert.Empty(m.CC)
	WithCC("cc0@bar.com", "cc1@bar.com")(&m)
	assert.Equal([]string{"cc0@bar.com", "cc1@bar.com"}, m.CC)

	assert.Empty(m.BCC)
	WithBCC("bcc0@bar.com", "bcc1@bar.com")(&m)
	assert.Equal([]string{"bcc0@bar.com", "bcc1@bar.com"}, m.BCC)

	assert.Empty(m.Subject)
	WithSubject("subject0")(&m)
	assert.Equal("subject0", m.Subject)

	assert.Empty(m.TextBody)
	WithTextBody("text body etc.")(&m)
	assert.Equal("text body etc.", m.TextBody)

	assert.Empty(m.HTMLBody)
	WithHTMLBody("html body etc.")(&m)
	assert.Equal("html body etc.", m.HTMLBody)

	assert.Equal("foo@bar.com", m.From)
	assert.Equal([]string{"buzz@bar.com", "fuzz@bar.com"}, m.To)
	assert.Equal([]string{"cc0@bar.com", "cc1@bar.com"}, m.CC)
	assert.Equal([]string{"bcc0@bar.com", "bcc1@bar.com"}, m.BCC)
	assert.Equal("text body etc.", m.TextBody)
	assert.Equal("html body etc.", m.HTMLBody)
}

func TestApplyMessageOptions(t *testing.T) {
	assert := assert.New(t)

	m := ApplyMessageOptions(Message{}, WithFrom("foo@bar.com"),
		WithTo("buzz@bar.com", "fuzz@bar.com"),
		WithCC("cc0@bar.com", "cc1@bar.com"),
		WithBCC("bcc0@bar.com", "bcc1@bar.com"),
		WithTextBody("text body etc."),
		WithHTMLBody("html body etc."))

	assert.Equal("foo@bar.com", m.From)
	assert.Equal([]string{"buzz@bar.com", "fuzz@bar.com"}, m.To)
	assert.Equal([]string{"cc0@bar.com", "cc1@bar.com"}, m.CC)
	assert.Equal([]string{"bcc0@bar.com", "bcc1@bar.com"}, m.BCC)
	assert.Equal("text body etc.", m.TextBody)
	assert.Equal("html body etc.", m.HTMLBody)
}
