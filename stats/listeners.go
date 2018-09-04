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
		var route string
		if len(wre.Route()) > 0 {
			route = fmt.Sprintf("route:%s", wre.Route())
		} else {
			route = "route:not_found"
		}

		method := fmt.Sprintf("method:%s", wre.Request().Method)

		tags := []string{
			route, method,
		}

		stats.Increment("http.request", tags...)
		if wre.StatusCode() >= http.StatusInternalServerError {
			stats.Increment("http.request.5xx", tags...)
		} else if wre.StatusCode() >= http.StatusBadRequest {
			stats.Increment("http.request.4xx", tags...)
		} else if wre.StatusCode() >= http.StatusMultipleChoices {
			stats.Increment("http.request.3xx", tags...)
		} else if wre.StatusCode() >= http.StatusOK {
			stats.Increment("http.request.2xx", tags...)
		}

		stats.Gauge("http.request.elapsed", util.Time.Millis(wre.Elapsed()), tags...)
		stats.Histogram("http.request.elapsed", util.Time.Millis(wre.Elapsed()), tags...)
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
		stats.Histogram("db.query.elapsed", util.Time.Millis(qe.Elapsed()), labelTag)
	}))
}

// AddErrorListeners adds error listeners.
func AddErrorListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.Warning, "stats", logger.NewErrorEventListener(func(qe *logger.ErrorEvent) {
		stats.Increment("warning")
	}))
	log.Listen(logger.Error, "stats", logger.NewErrorEventListener(func(qe *logger.ErrorEvent) {
		stats.Increment("error")
	}))
	log.Listen(logger.Fatal, "stats", logger.NewErrorEventListener(func(qe *logger.ErrorEvent) {
		stats.Increment("fatal")
	}))
}
