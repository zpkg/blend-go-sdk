package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestHealthz(t *testing.T) {
	t.Skip()
	assert := assert.New(t)

	appLog := logger.New().WithFlags(logger.AllFlags())
	defer appLog.Close()

	app := New().WithBindAddr("127.0.0.1:0").WithLogger(appLog).WithConfig(MustNewConfigFromEnv())
	defer app.Shutdown()

	appStarted := make(chan struct{})
	appLog.Listen(AppStartComplete, "default", NewAppEventListener(func(aes *AppEvent) {
		close(appStarted)
	}))

	hzLog := logger.New().WithFlags(logger.AllFlags())
	defer hzLog.Close()

	hz := NewHealthz(app).WithLogger(hzLog).WithGracePeriod(0)
	defer hz.Shutdown()
	hz.WithDefaultHeader("key", "secure")
	assert.NotEmpty(hz.DefaultHeaders())

	assert.NotNil(hz.Hosted())
	assert.False(app.Latch().IsRunning())

	go hz.Start()
	<-hz.hosted.NotifyStarted()
	<-hz.self.NotifyStarted()

	assert.True(hz.hosted.IsRunning())
	assert.True(hz.self.Latch().IsRunning())

	assert.NotNil(hz.self.Listener())

	healthzRes, err := http.Get("http://" + hz.self.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusOK, healthzRes.StatusCode)
	assert.Equal("secure", healthzRes.Header.Get("key"))

	app.Shutdown()
	<-app.NotifyShutdown()

	healthzRes, err = http.Get("http://" + hz.self.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, healthzRes.StatusCode)
}

func TestHealthzShutdown(t *testing.T) {
	t.Skip()
	assert := assert.New(t)

	appLog := logger.New().WithFlags(logger.AllFlags())
	defer appLog.Close()

	app := New().WithBindAddr("127.0.0.1:0").WithLogger(appLog).WithConfig(MustNewConfigFromEnv())
	defer app.Shutdown()

	appStarted := make(chan struct{})
	appLog.Listen(AppStartComplete, "default", NewAppEventListener(func(aes *AppEvent) {
		close(appStarted)
	}))

	hzLog := logger.New().WithFlags(logger.AllFlags())
	defer hzLog.Close()

	hz := NewHealthz(app).
		WithLogger(hzLog).
		WithGracePeriod(30 * time.Second).
		WithFailureThreshold(3)

	assert.NotNil(hz.Hosted())
	assert.False(app.Latch().IsRunning())

	go hz.Start()
	<-hz.hosted.NotifyStarted()
	<-hz.self.NotifyStarted()
	<-hz.NotifyStarted()

	assert.True(hz.IsRunning())

	// shutdown the server
	go hz.Shutdown()
	<-hz.NotifyShuttingDown()

	assert.True(hz.latch.IsStopping())

	res, err := http.Get("http://" + hz.self.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, res.StatusCode)

	assert.True(hz.IsRunning())

	res, err = http.Get("http://" + hz.self.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, res.StatusCode)

	assert.True(hz.IsRunning())

	res, err = http.Get("http://" + hz.self.Listener().Addr().String() + "/healthz")
	assert.Nil(err)
	assert.Equal(http.StatusServiceUnavailable, res.StatusCode)

	<-hz.NotifyShutdown()
	<-hz.self.NotifyShutdown()
	<-hz.hosted.NotifyShutdown()

	assert.False(hz.IsRunning())
	assert.False(hz.self.IsRunning())
	assert.False(hz.hosted.IsRunning())
}

func TestHealthzProperties(t *testing.T) {
	assert := assert.New(t)

	hz := NewHealthz(nil)
	assert.True(hz.RecoverPanics())
	hz.WithRecoverPanics(false)
	assert.False(hz.RecoverPanics())

	assert.Nil(hz.Logger())
	hz.WithLogger(logger.None())
	assert.NotNil(hz.Logger())
}
