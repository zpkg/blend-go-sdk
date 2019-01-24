package diagnostics

import "net/http"

// Notifier is a diagnostics notifier.
type Notifier interface {
	Notify(err interface{}) error
	NotifyWithRequest(err interface{}, req *http.Request) error
}
