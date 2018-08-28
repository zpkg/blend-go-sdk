package request

import (
	"github.com/blend/go-sdk/logger"
)

// NewManager creates a new manager.
func NewManager() *Manager {
	return &Manager{}
}

// Manager is a helper to create requests with common metadata.
// It is generally for creating requests to *any* host.
type Manager struct {
	Log                    *logger.Logger
	MockedResponseProvider MockedResponseProvider
	OnRequest              Handler
	OnResponse             ResponseHandler
	Tracer                 Tracer
}

// WithLogger sets the logger.
func (m *Manager) WithLogger(log *logger.Logger) *Manager {
	m.Log = log
	return m
}

// WithMockedResponseProvider sets the mocked response provider.
func (m *Manager) WithMockedResponseProvider(mrp MockedResponseProvider) *Manager {
	m.MockedResponseProvider = mrp
	return m
}

// WithOnRequest sets the on request handler..
func (m *Manager) WithOnRequest(handler Handler) *Manager {
	m.OnRequest = handler
	return m
}

// WithOnResponse sets the on response handler.
func (m *Manager) WithOnResponse(handler ResponseHandler) *Manager {
	m.OnResponse = handler
	return m
}

// WithTracer sets the tracer.
func (m *Manager) WithTracer(tracer Tracer) *Manager {
	m.Tracer = tracer
	return m
}

// Create creates a new request.
func (m Manager) Create() *Request {
	return New().
		WithLogger(m.Log).
		WithMockProvider(m.MockedResponseProvider).
		WithRequestHandler(m.OnRequest).
		WithResponseHandler(m.OnResponse).
		WithTracer(m.Tracer)
}

// Get returns a new get request for a given url.
func (m Manager) Get(url string) (*Request, error) {
	return m.Create().AsGet().WithRawURL(url)
}

// Post returns a new post request for a given url.
func (m Manager) Post(url string, body []byte) (*Request, error) {
	return m.Create().AsPost().WithPostBody(body).WithRawURL(url)
}
