package main

import (
	"context"
	"flag"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
)

var exec = flag.String("exec", "", "The command to execute")
var name = flag.String("exec", stringutil.Letters.Random(8), "The command to execute")
var schedule = flag.String("schedule", "* */5 * * *", "The job schedule")
var config = flag.String("config", "config.yml", "The job manager config")

func main() {
	flag.Parse()

	schedule, err := cron.ParseString(*schedule)
	if err != nil {
		logger.FatalExit(err)
	}

	jm := cron.New()
	jm.LoadJob(&Job{
		schedule: schedule,
		name:     *name,
		exec:     *exec,
	})

	if err := graceful.Shutdown(jm); err != nil {
		logger.FatalExit(err)
	}
}

// Job is the main job body.
type Job struct {
	schedule *cron.StringSchedule
	config   *jobkit.Config
	name     string
	exec     string
}

// Name returns the job name.
func (job Job) Name() string {
	return job.name
}

// Schedule returns the job schedule.
func (job Job) Schedule() cron.Schedule {
	return job.schedule
}

// Execute is the job body.
func (job Job) Execute(ctx context.Context) error {
	return sh.Exec(job.exec)
}
