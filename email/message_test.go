package email

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestMessageValidate(t *testing.T) {
	assert := assert.New(t)

	assert.True(exception.Is(ErrMessageFieldUnset, Message{}.Validate()))
	assert.True(exception.Is(ErrMessageFieldUnset, Message{
		From: "foo@bar.com",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From: "foo\r@bar.com",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From: "foo\n@bar.com",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From: "foo\r\n@bar.com",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com", "bad\n@bar.com"},
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com"},
		CC:   []string{"bad\n@bar.com"},
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com"},
		CC:   []string{"ok@bar.com"},
		BCC:  []string{"bad\n@bar.com"},
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From:    "foo@bar.com",
		To:      []string{"moo@bar.com"},
		Subject: "this is \n bad",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From:    "foo@bar.com",
		To:      []string{"moo@bar.com"},
		Subject: "this is \r bad",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldNewlines, Message{
		From:    "foo@bar.com",
		To:      []string{"moo@bar.com"},
		Subject: "this is \n\r bad",
	}.Validate()))
	assert.True(exception.Is(ErrMessageFieldUnset, Message{
		From: "foo@bar.com",
		To:   []string{"moo@bar.com"},
	}.Validate()))

	assert.Nil(Message{
		From:     "foo@bar.com",
		To:       []string{"moo@bar.com"},
		TextBody: "stuff",
	}.Validate())
}
