package diagnostics

import "net/http"

var (
	_ Notifier = (*MockNotifier)(nil)
)

// MockNotifier is a testing notifier.
type MockNotifier chan MockNotification

// Notify notifies with an error.
func (mn MockNotifier) Notify(err interface{}) error {
	mn <- MockNotification{Err: err}
	return nil
}

// NotifyWithRequest notifies with a request.
func (mn MockNotifier) NotifyWithRequest(err interface{}, req *http.Request) error {
	mn <- MockNotification{Err: err, Req: req}
	return nil
}

// MockNotification is a mocked notification.
type MockNotification struct {
	Err interface{}
	Req *http.Request
}
