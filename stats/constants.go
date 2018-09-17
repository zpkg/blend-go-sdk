package stats

import "github.com/blend/go-sdk/logger"

// MetricNames are names we use when sending data to the collectors.
const (
	MetricNameHTTPRequest        string = string(logger.HTTPRequest)
	MetricNameHTTPRequestElapsed string = MetricNameHTTPRequest + ".elapsed"
	MetricNameDBQuery            string = string(logger.Query)
	MetricNameDBQueryElapsed     string = MetricNameDBQuery + ".elapsed"

	MetricNameError string = string(logger.Error)

	TagService   string = "service"
	TagEnv       string = "env"
	TagContainer string = "container"

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
