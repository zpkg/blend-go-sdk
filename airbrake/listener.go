package airbrake

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
)

const (
	// ListenerAirbrake is the airbrake listener name.
	ListenerAirbrake = "airbrake"
)

// AddListeners adds airbrake listeners.
func AddListeners(log logger.Listenable, cfg *Config) {
	if log == nil || cfg == nil || cfg.IsZero() {
		return
	}
	client := New(cfg)
	listener := logger.NewErrorEventListener(func(ee *logger.ErrorEvent) {
		if req, ok := ee.State().(*http.Request); ok {
			client.NotifyWithRequest(ee.Err(), req)
		} else {
			client.Notify(ee.Err())
		}
	})
	log.Listen(logger.Error, ListenerAirbrake, listener)
	log.Listen(logger.Fatal, ListenerAirbrake, listener)
}
