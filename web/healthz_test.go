package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestHealthz(t *testing.T) {
	assert := assert.New(t)

	appLog := logger.New().WithFlags(logger.AllFlags())
	defer appLog.Close()

	app := New().WithBindAddr("127.0.0.1:0").WithLogger(appLog)
	defer app.Shutdown()

	appStarted := make(chan struct{})
	appLog.Listen(AppStartComplete, "default", NewAppEventListener(func(aes *AppEvent) {
		close(appStarted)
	}))

	hzLog := logger.New().WithFlags(logger.AllFlags())
	defer hzLog.Close()

	hz := NewHealthz(app).WithLogger(hzLog)
	hz.WithDefaultHeader("key", "secure")
	assert.NotEmpty(hz.DefaultHeaders())
	hzApp := New().WithBindAddr("127.0.0.1:0").WithHandler(hz)

	assert.NotNil(hz.App())
	assert.False(app.Latch().IsRunning())

	go app.Start()
	go hzApp.Start()
	<-app.NotifyStarted()
	<-hzApp.NotifyStarted()

	assert.True(app.Latch().IsRunning())
	assert.True(hz.App().Latch().IsRunning())
	assert.NotNil(hzApp.Listener())

	healthzRes, err := http.Get("http://" + hzApp.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusOK, healthzRes.StatusCode)
	assert.Equal("secure", healthzRes.Header.Get("key"))

	app.Shutdown()

	healthzRes, err = http.Get("http://" + hzApp.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, healthzRes.StatusCode)

	varzRes, err := http.Get("http://" + hzApp.Listener().Addr().String() + "/varz")
	assert.Nil(err)
	assert.Equal(http.StatusOK, varzRes.StatusCode)

	notfoundRes, err := http.Get("http://" + hzApp.Listener().Addr().String() + "/adfasdfa")
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, notfoundRes.StatusCode)
}

func TestHealthzHTTPResponseListener(t *testing.T) {
	assert := assert.New(t)

	hz := NewHealthz(nil)
	hz.httpResponseListener((&logger.HTTPResponseEvent{}).WithStatusCode(http.StatusOK))
	assert.Equal(1, hz.vars.Get(VarzRequests))
	assert.Equal(1, hz.vars.Get(VarzRequests2xx))
	assert.Equal(0, hz.vars.Get(VarzRequests3xx))
	assert.Equal(0, hz.vars.Get(VarzRequests4xx))
	assert.Equal(0, hz.vars.Get(VarzRequests5xx))

	hz.httpResponseListener((&logger.HTTPResponseEvent{}).WithStatusCode(http.StatusBadRequest))
	assert.Equal(2, hz.vars.Get(VarzRequests))
	assert.Equal(1, hz.vars.Get(VarzRequests2xx))
	assert.Equal(0, hz.vars.Get(VarzRequests3xx))
	assert.Equal(1, hz.vars.Get(VarzRequests4xx))
	assert.Equal(0, hz.vars.Get(VarzRequests5xx))

	hz.httpResponseListener((&logger.HTTPResponseEvent{}).WithStatusCode(http.StatusInternalServerError))
	assert.Equal(3, hz.vars.Get(VarzRequests))
	assert.Equal(1, hz.vars.Get(VarzRequests2xx))
	assert.Equal(0, hz.vars.Get(VarzRequests3xx))
	assert.Equal(1, hz.vars.Get(VarzRequests4xx))
	assert.Equal(1, hz.vars.Get(VarzRequests5xx))

	hz.httpResponseListener((&logger.HTTPResponseEvent{}).WithStatusCode(http.StatusMovedPermanently))
	assert.Equal(4, hz.vars.Get(VarzRequests))
	assert.Equal(1, hz.vars.Get(VarzRequests2xx))
	assert.Equal(1, hz.vars.Get(VarzRequests3xx))
	assert.Equal(1, hz.vars.Get(VarzRequests4xx))
	assert.Equal(1, hz.vars.Get(VarzRequests5xx))
}

func TestHealthzErrorListener(t *testing.T) {
	assert := assert.New(t)

	hz := NewHealthz(nil)
	hz.errorListener(logger.Errorf(logger.Error, ""))
	assert.Equal(1, hz.vars.Get(VarzErrors))
	assert.Equal(0, hz.vars.Get(VarzFatals))

	hz.errorListener(logger.Errorf(logger.Fatal, ""))
	assert.Equal(1, hz.vars.Get(VarzErrors))
	assert.Equal(1, hz.vars.Get(VarzFatals))
}

func TestHealthzProperties(t *testing.T) {
	assert := assert.New(t)

	hz := NewHealthz(nil)
	assert.False(hz.RecoverPanics())
	hz.WithRecoverPanics(true)
	assert.True(hz.RecoverPanics())

	assert.Nil(hz.Logger())
	hz.WithLogger(logger.None())
	assert.NotNil(hz.Logger())
}

func TestHealthzEnsureListeners(t *testing.T) {
	assert := assert.New(t)

	app := New().WithLogger(logger.None())
	hz := NewHealthz(app)
	hz.ensureListeners()

	started, _ := hz.Vars().Get(VarzStarted).(time.Time)
	assert.False(started.IsZero())
	assert.True(app.Logger().HasListener(logger.HTTPResponse, ListenerHealthz))
	assert.True(app.Logger().HasListener(logger.Error, ListenerHealthz))
	assert.True(app.Logger().HasListener(logger.Fatal, ListenerHealthz))

	// shouldn't do anything
	hz.ensureListeners()
}
