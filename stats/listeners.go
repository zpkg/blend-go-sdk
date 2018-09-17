package stats

import (
	"strconv"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
)

// MetricNames are names we use when sending data to the collectors.
const (
	MetricNameHTTPRequest        string = string(logger.HTTPRequest)
	MetricNameHTTPRequestElapsed string = MetricNameHTTPRequest + ".elapsed"
	MetricNameDBQuery            string = string(logger.Query)
	MetricNameDBQueryElapsed     string = MetricNameDBQuery + ".elapsed"

	MetricNameError string = string(logger.Error)

	TagRoute  string = "route"
	TagMethod string = "method"
	TagStatus string = "status"

	TagQuery    string = "query"
	TagEngine   string = "engine"
	TagDatabase string = "database"

	TagSeverity string = "severity"
	TagError    string = "error"
	TagClass    string = "class"

	RouteNotFound string = "not_found"

	ListenerNameStats string = "stats"
)

// Tag creates a new tag.
func Tag(key, value string) string {
	return key + ":" + value
}

// AddWebListeners adds web listeners.
func AddWebListeners(log *logger.Logger, stats Collector) {
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

		elapsed := util.Time.Millis(wre.Elapsed())
		stats.Increment(MetricNameHTTPRequest, tags...)
		stats.Gauge(MetricNameHTTPRequestElapsed, elapsed, tags...)
		stats.Histogram(MetricNameHTTPRequestElapsed, elapsed, tags...)
	}))
}

// AddQueryListeners adds db listeners.
func AddQueryListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.Query, ListenerNameStats, logger.NewQueryEventListener(func(qe *logger.QueryEvent) {
		engine := Tag(TagEngine, qe.Engine())
		database := Tag(TagDatabase, qe.Database())

		tags := []string{
			engine, database,
		}
		if len(qe.QueryLabel()) > 0 {
			tags = append(tags, Tag(TagQuery, qe.QueryLabel()))
		}
		if qe.Err() != nil {
			if ex := exception.As(qe.Err()); ex != nil && ex.Class() != nil {
				tags = append(tags, Tag(TagClass, ex.Class().Error()))
			}
			tags = append(tags, TagError)
		}

		stats.Increment(MetricNameDBQuery, tags...)
		stats.Gauge(MetricNameDBQueryElapsed, util.Time.Millis(qe.Elapsed()), tags...)
		stats.Histogram(MetricNameDBQueryElapsed, util.Time.Millis(qe.Elapsed()), tags...)
	}))
}

// AddErrorListeners adds error listeners.
func AddErrorListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	listener := logger.NewErrorEventListener(func(ee *logger.ErrorEvent) {
		stats.Increment(MetricNameError,
			Tag(TagSeverity, string(ee.Flag())),
			Tag(TagClass, exception.ErrClass(ee.Err())),
		)
	})
	log.Listen(logger.Warning, ListenerNameStats, listener)
	log.Listen(logger.Error, ListenerNameStats, listener)
	log.Listen(logger.Fatal, ListenerNameStats, listener)
}
