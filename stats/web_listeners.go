package stats

import (
	"strconv"

	"github.com/blend/go-sdk/logger"
)

// AddWebListeners adds web listeners.
func AddWebListeners(log logger.Log, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.HTTPResponse, ListenerNameStats, logger.NewHTTPResponseEventListener(func(wre *logger.HTTPResponseEvent) {
		var route string
		if len(wre.Route()) > 0 {
			route = Tag(TagRoute, wre.Route())
		} else {
			route = Tag(TagRoute, RouteNotFound)
		}

		method := Tag(TagMethod, wre.Request().Method)
		status := Tag(TagStatus, strconv.Itoa(wre.StatusCode()))
		tags := []string{
			route, method, status,
		}

		stats.Increment(MetricNameHTTPRequest, tags...)
		stats.TimeInMilliseconds(MetricNameHTTPRequestElapsed, wre.Elapsed(), tags...)
	}))
}
