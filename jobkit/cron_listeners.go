package jobkit

import (
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

// ListenerNames
const (
	ListenerNameJobKit = "jobkit"
)

// AddCronListeners adds listeners for cron tasks.
func AddCronListeners(log *logger.Logger) {
	if log == nil {
		return
	}

	log.Listen(cron.FlagStarted, ListenerNameJobUtil, cron.NewEventListener(func(ce *cron.Event) {

	}))
}
