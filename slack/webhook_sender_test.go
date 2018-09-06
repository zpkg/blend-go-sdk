package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWebhookSender(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	config := &Config{
		Webhook: ts.URL,
	}

	sender := NewWebhookSender(config)
	assert.Nil(sender.Send(&Message{
		Text: "this is only a test",
	}))
}

func TestWebhookSenderStatusCode(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	config := &Config{
		Webhook: ts.URL,
	}

	sender := NewWebhookSender(config)
	assert.NotNil(sender.Send(&Message{
		Text: "this is only a test",
	}))
}
