package slack

import (
	"context"
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

	sender := New(config)
	assert.Nil(sender.Send(context.TODO(), Message{
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

	sender := New(config)
	assert.NotNil(sender.Send(context.TODO(), Message{
		Text: "this is only a test",
	}))
}

func TestWebhookSenderApplyDefaults(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Webhook:  "http://foo.com",
		Channel:  "#bot-test",
		Username: "default-test",
	}

	sender := New(config)
	updated := sender.ApplyDefaults(Message{
		Text: "this is only a test",
	})

	assert.Equal("this is only a test", updated.Text)
	assert.Equal("#bot-test", updated.Channel)
	assert.Equal("default-test", updated.Username)
}
