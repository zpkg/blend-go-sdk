package stats

import (
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
)

// AddWebListeners adds web listeners.
func AddWebListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.HTTPResponse, "stats", logger.NewHTTPResponseEventListener(func(wre *logger.HTTPResponseEvent) {
		var path string
		if len(wre.Route()) > 0 {
			path = fmt.Sprintf("path:%s", wre.Route())
		} else {
			path = "path:not_found"
		}

		stats.Increment("request", path)
		if wre.StatusCode() >= http.StatusInternalServerError {
			stats.Increment("request.5xx", path)
		} else if wre.StatusCode() >= http.StatusBadRequest {
			stats.Increment("request.4xx", path)
		} else if wre.StatusCode() >= http.StatusMultipleChoices {
			stats.Increment("request.3xx", path)
		} else if wre.StatusCode() >= http.StatusOK {
			stats.Increment("request.2xx", path)
		}

		stats.Gauge("request.elapsed", util.Time.Millis(wre.Elapsed()), path)
		stats.Histogram("request.elapsed", util.Time.Millis(wre.Elapsed()), path)
	}))
}

// AddQueryListeners adds db listeners.
func AddQueryListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.Query, "stats", logger.NewQueryEventListener(func(qe *logger.QueryEvent) {
		if len(qe.QueryLabel()) == 0 {
			return
		}

		labelTag := fmt.Sprintf("query:%s", qe.QueryLabel())

		stats.Increment("db.query", labelTag)
		if qe.Err() != nil {
			stats.Increment("db.query.error", labelTag)
		}

		stats.Gauge("db.query.elapsed", util.Time.Millis(qe.Elapsed()), labelTag)
		stats.Histogram("request.elapsed", util.Time.Millis(qe.Elapsed()), labelTag)
	}))
}
