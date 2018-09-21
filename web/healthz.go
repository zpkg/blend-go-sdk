package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// Verify interfaces.
var (
	_ Shutdowner = (*Healthz)(nil)
)

const (
	// VarzStarted is a common variable.
	VarzStarted = "startedUTC"
	// VarzRequests is a common variable.
	VarzRequests = "http_requests"
	// VarzRequests2xx is a common variable.
	VarzRequests2xx = "http_requests2xx"
	// VarzRequests3xx is a common variable.
	VarzRequests3xx = "http_requests3xx"
	// VarzRequests4xx is a common variable.
	VarzRequests4xx = "http_requests4xx"
	// VarzRequests5xx is a common variable.
	VarzRequests5xx = "http_requests5xx"
	// VarzErrors is a common variable.
	VarzErrors = "errors_total"
	// VarzFatals is a common variable.
	VarzFatals = "fatals_total"

	// ListenerHealthz is the uid of the healthz logger listeners.
	ListenerHealthz = "healthz"

	// ErrHealthzAppUnset is a common error.
	ErrHealthzAppUnset exception.Class = "healthz app unset"
)

// NewHealthz returns a new healthz.
func NewHealthz(hosted Shutdowner) *Healthz {
	return &Healthz{
		hosted:         hosted,
		bindAddr:       DefaultHealthzBindAddr,
		gracePeriod:    DefaultShutdownGracePeriod,
		readTimeout:    DefaultReadTimeout,
		writeTimeout:   DefaultWriteTimeout,
		latch:          &async.Latch{},
		defaultHeaders: map[string]string{},
		recoverPanics:  true,
	}
}

// Healthz is a sentinel / healthcheck sidecar that can run on a different
// port to the main app.
/*
It typically implements the following routes:

	/healthz - overall health endpoint, 200 on healthy, 5xx on not.
				should be used as a kubernetes readiness probe.
	/debug/vars - `pkg/expvar` output.
*/
type Healthz struct {
	self           *App
	hosted         Shutdowner
	cfg            *HealthzConfig
	bindAddr       string
	log            *logger.Logger
	latch          *async.Latch
	defaultHeaders map[string]string
	recoverPanics  bool

	maxHeaderBytes    int
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	gracePeriod       time.Duration

	failureThreshold int
	failures         int32
}

// WithConfig sets the healthz config and relevant properties.
func (hz *Healthz) WithConfig(cfg *HealthzConfig) *Healthz {
	hz.cfg = cfg
	hz.WithBindAddr(cfg.GetBindAddr())
	hz.WithGracePeriod(cfg.GetGracePeriod())
	hz.WithFailureThreshold(cfg.GetFailureThreshold())
	hz.WithRecoverPanics(cfg.GetRecoverPanics())
	hz.WithMaxHeaderBytes(cfg.GetMaxHeaderBytes())
	hz.WithReadHeaderTimeout(cfg.GetReadHeaderTimeout())
	hz.WithReadTimeout(cfg.GetReadTimeout())
	hz.WithWriteTimeout(cfg.GetWriteTimeout())
	hz.WithIdleTimeout(cfg.GetIdleTimeout())
	return hz
}

// Config returns the healthz config.
func (hz *Healthz) Config() *HealthzConfig {
	return hz.cfg
}

// WithBindAddr sets the bind address.
func (hz *Healthz) WithBindAddr(bindAddr string) *Healthz {
	hz.bindAddr = bindAddr
	return hz
}

// BindAddr returns the bind address.
func (hz *Healthz) BindAddr() string {
	return hz.bindAddr
}

// WithGracePeriod sets the grace period seconds
func (hz *Healthz) WithGracePeriod(gracePeriod time.Duration) *Healthz {
	hz.gracePeriod = gracePeriod
	return hz
}

// GracePeriod returns the grace period in seconds
func (hz *Healthz) GracePeriod() time.Duration {
	return hz.gracePeriod
}

// WithFailureThreshold sets the failure threshold.
func (hz *Healthz) WithFailureThreshold(failureThreshold int) *Healthz {
	hz.failureThreshold = failureThreshold
	return hz
}

// FailureThreshold returns the failure threshold.
func (hz *Healthz) FailureThreshold() int {
	return hz.failureThreshold
}

// Hosted returns the underlying app.
func (hz *Healthz) Hosted() Shutdowner {
	return hz.hosted
}

// RecoverPanics returns if the app recovers panics.
func (hz *Healthz) RecoverPanics() bool {
	return hz.recoverPanics
}

// WithRecoverPanics sets if the app should recover panics.
func (hz *Healthz) WithRecoverPanics(value bool) *Healthz {
	hz.recoverPanics = value
	return hz
}

// Logger returns the diagnostics agent for the app.
func (hz *Healthz) Logger() *logger.Logger {
	return hz.log
}

// WithLogger sets the app logger agent and returns a reference to the app.
// It also sets underlying loggers in any child resources like providers and the auth manager.
func (hz *Healthz) WithLogger(log *logger.Logger) *Healthz {
	hz.log = log
	return hz
}

// WithDefaultHeader adds a default header.
func (hz *Healthz) WithDefaultHeader(key, value string) *Healthz {
	hz.defaultHeaders[key] = value
	return hz
}

// DefaultHeaders returns the default headers.
func (hz *Healthz) DefaultHeaders() map[string]string {
	return hz.defaultHeaders
}

// WithMaxHeaderBytes sets the max header bytes value and returns a reference.
func (hz *Healthz) WithMaxHeaderBytes(byteCount int) *Healthz {
	hz.maxHeaderBytes = byteCount
	return hz
}

// MaxHeaderBytes returns the app max header bytes.
func (hz *Healthz) MaxHeaderBytes() int {
	return hz.maxHeaderBytes
}

// WithReadHeaderTimeout returns the read header timeout for the server.
func (hz *Healthz) WithReadHeaderTimeout(timeout time.Duration) *Healthz {
	hz.readHeaderTimeout = timeout
	return hz
}

// ReadHeaderTimeout returns the read header timeout for the server.
func (hz *Healthz) ReadHeaderTimeout() time.Duration {
	return hz.readHeaderTimeout
}

// WithReadTimeout sets the read timeout for the server and returns a reference to the app for building apps with a fluent api.
func (hz *Healthz) WithReadTimeout(timeout time.Duration) *Healthz {
	hz.readTimeout = timeout
	return hz
}

// ReadTimeout returns the read timeout for the server.
func (hz *Healthz) ReadTimeout() time.Duration {
	return hz.readTimeout
}

// WithIdleTimeout sets the idle timeout.
func (hz *Healthz) WithIdleTimeout(timeout time.Duration) *Healthz {
	hz.idleTimeout = timeout
	return hz
}

// IdleTimeout is the time before we close a connection.
func (hz *Healthz) IdleTimeout() time.Duration {
	return hz.idleTimeout
}

// WithWriteTimeout sets the write timeout for the server and returns a reference to the app for building apps with a fluent api.
func (hz *Healthz) WithWriteTimeout(timeout time.Duration) *Healthz {
	hz.writeTimeout = timeout
	return hz
}

// WriteTimeout returns the write timeout for the server.
func (hz *Healthz) WriteTimeout() time.Duration {
	return hz.writeTimeout
}

// Start implements shutdowner.
func (hz *Healthz) Start() error {
	hz.latch.Starting()
	hz.self = New().
		WithHandler(hz).
		WithBindAddr(hz.bindAddr).
		WithLogger(hz.log)

	hz.latch.Started()
	return async.RunToError(hz.self.Start, hz.hosted.Start)
}

// Shutdown implements shutdowner.
func (hz *Healthz) Shutdown() error {
	// set the next call to `/healtz` to
	// finish the shutdown
	hz.latch.Stopping()
	defer func() { hz.latch.Stopped() }()

	context, cancel := context.WithTimeout(context.Background(), hz.GracePeriod())
	defer cancel()

	if hz.log != nil {
		hz.log.Infof("healthz is shutting down with (%s) grace period", hz.GracePeriod())
	}

	select {
	// if the hosted app crashes
	case <-hz.hosted.NotifyShutdown():
		return hz.self.Shutdown()
	// if the shutdown grace period expires
	case <-context.Done():
		if hz.log != nil {
			hz.log.Warningf("healthz shutdown grace period has expired")
		}
		return hz.shutdownServers()
	// if we've received a final /healthz request
	case <-hz.latch.NotifyStopped():
		return hz.shutdownServers()
	}
}

// IsRunning returns if the healthz server is running.
func (hz *Healthz) IsRunning() bool {
	return hz.self.IsRunning()
}

// NotifyStarted returns the notify started signal.
func (hz *Healthz) NotifyStarted() <-chan struct{} {
	return hz.self.NotifyStarted()
}

// NotifyShuttingDown returns the notify shutdown signal.
func (hz *Healthz) NotifyShuttingDown() <-chan struct{} {
	return hz.latch.NotifyStopping()
}

// NotifyShutdown returns the notify shutdown signal.
func (hz *Healthz) NotifyShutdown() <-chan struct{} {
	return hz.latch.NotifyStopped()
}

func (hz *Healthz) shutdownServers() error {
	return async.RunToError(hz.hosted.Shutdown, hz.self.Shutdown)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (hz *Healthz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hz.recoverPanics {
		defer hz.recover(w, r)
	}

	start := time.Now()
	route := strings.ToLower(r.URL.Path)

	res := NewRawResponseWriter(w)
	res.Header().Set(HeaderContentEncoding, ContentEncodingIdentity)

	if hz.log != nil {
		hz.log.Trigger(logger.NewHTTPRequestEvent(r).WithRoute(route))
		defer func() {
			hz.log.Trigger(logger.NewHTTPResponseEvent(r).
				WithStatusCode(res.StatusCode()).
				WithElapsed(time.Since(start)).
				WithContentLength(res.ContentLength()),
			)
		}()
	}

	if len(hz.defaultHeaders) > 0 {
		for key, value := range hz.defaultHeaders {
			res.Header().Set(key, value)
		}
	}

	switch route {
	case "/healthz":
		hz.healthzHandler(res, r)
	default:
		http.NotFound(res, r)
	}

	if err := res.Close(); err != nil && err != http.ErrBodyNotAllowed && hz.log != nil {
		hz.log.Error(err)
	}
}

func (hz *Healthz) recover(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		if hz.log != nil {
			hz.log.Fatalf("%v", rcv)
		}

		http.Error(w, fmt.Sprintf("%v", rcv), http.StatusInternalServerError)
		return
	}
}

func (hz *Healthz) requiredFailures() int32 {
	if hz.cfg != nil {
		return int32(hz.cfg.GetFailureThreshold())
	}
	return DefaultHealthzFailureThreshold
}

// currentFailures returns the current failures.
func (hz *Healthz) currentFailures() (output int32) {
	output = atomic.LoadInt32(&hz.failures)
	return
}

// incrementFailures increments failures.
func (hz *Healthz) incrementFailures() {
	atomic.AddInt32(&hz.failures, 1)
}

// resetFailures resets the failures count.
func (hz *Healthz) resetFailures() {
	atomic.StoreInt32(&hz.failures, 0)
}

func (hz *Healthz) healthzHandler(w ResponseWriter, r *http.Request) {
	if hz.latch.IsStopping() {
		hz.incrementFailures()

		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set(HeaderContentType, ContentTypeText)
		fmt.Fprintf(w, "Shutting down.\n")
		if hz.log != nil {
			hz.log.Debugf("healthz received probe while in process of shutdown")
		}

		// handle max fails ...
		if hz.currentFailures() >= int32(hz.FailureThreshold()) {
			hz.latch.Stopped()
		}
	} else if hz.hosted.IsRunning() {
		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, ContentTypeText)
		fmt.Fprintf(w, "OK!\n")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set(HeaderContentType, ContentTypeText)
		fmt.Fprintf(w, "Failure!\n")
	}
	return
}
