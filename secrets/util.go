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

// MetaProvider is a service path meta provider.
type MetaProvider interface {
	GetServiceName() string
	GetServiceEnv() string
}

// ServicePath is the service key path.
func ServicePath(metaProvider MetaProvider) string {
	name := metaProvider.GetServiceName()
	environment := metaProvider.GetServiceEnv()
	return fmt.Sprintf("/services/%s/%s", environment, name)
}

// ServiceConfigPath returns the service config path.
func ServiceConfigPath(metaProvider MetaProvider) string {
	return fmt.Sprintf("%s/config", ServicePath(metaProvider))
}
