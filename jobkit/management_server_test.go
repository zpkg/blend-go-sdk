package jobkit

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/web"
)

func TestManagementServer(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()

	jm.LoadJob(cron.NewJob("test0"))
	jm.LoadJob(cron.NewJob("test1"))

	app := NewManagementServer(jm, &Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := app.Mock().Get("/").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	var jobs cron.Status
	meta, err = app.Mock().Get("/api/jobs").JSONWithMeta(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Len(jobs.Jobs, 2)
}

func TestManagementServerHealthz(t *testing.T) {
	assert := assert.New(t)

	jm := cron.New()
	jm.LoadJob(cron.NewJob("test0"))
	jm.LoadJob(cron.NewJob("test1"))
	jm.Start()

	app := NewManagementServer(jm, &Config{
		Web: web.Config{
			Port: 5000,
		},
	})

	meta, err := app.Mock().Get("/healthz").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	jm.Stop()

	meta, err = app.Mock().Get("/healthz").ExecuteWithMeta()
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, meta.StatusCode)
}
