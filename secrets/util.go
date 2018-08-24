package secrets

import (
	"fmt"
	"net/url"
)

// URL creates a new url.
func URL(format string, args ...interface{}) *url.URL {
	output, err := url.ParseRequestURI(fmt.Sprintf(format, args...))
	if err != nil {
		panic(err)
	}
	return output
}

// ServiceConfigPath returns the service config path.
func ServiceConfigPath(config Config) string {
	return fmt.Sprintf("%s/config", config.GetServicePath())
}
