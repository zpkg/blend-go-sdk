package web

import (
	"net/http"
	"testing"

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

	hz := NewHealthz(app).WithBindAddr("127.0.0.1:0").WithLogger(hzLog)
	hzServer := hz.Server()

	hzApp := New().WithServer(hzServer)

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
}
