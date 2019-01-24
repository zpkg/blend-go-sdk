package main

import (
	"context"
	"flag"
	"os"

	"github.com/blend/go-sdk/configutil"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
)

var name = flag.String("name", stringutil.Letters.Random(8), "The name of the job")
var exec = flag.String("exec", "", "The command to execute")
var schedule = flag.String("schedule", "*/1 * * * * * *", "The job schedule as a cron string (i.e. 7 space delimited components)")
var configPath = flag.String("config", "config.yml", "The job config path")

func main() {
	flag.Parse()

	schedule, err := cron.ParseString(*schedule)
	if err != nil {
		logger.FatalExit(err)
	}

	var config jobkit.Config
	if err := configutil.Read(&config, *configPath); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}

	log := logger.NewFromConfig(&config.Logger)
	log.WithEnabled(cron.FlagStarted, cron.FlagComplete, cron.FlagFixed, cron.FlagBroken, cron.FlagFailed, cron.FlagCancelled)

	command := *exec
	if command == "" {
		command, err = sh.ParseFlagsTrailer(os.Args...)
		if err != nil {
			logger.FatalExit(err)
		}
	}

	jm := cron.New().WithLogger(log)
	jm.LoadJob(&Job{
		schedule: schedule,
		name:     *name,
		exec:     command,
	})

	go func() {
		if err := graceful.Shutdown(jm); err != nil {
			logger.FatalExit(err)
		}
	}()
	ws := jobkit.NewManagementServer(jm, &config)
	ws.WithLogger(log)
	if err := graceful.Shutdown(ws); err != nil {
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
	return sh.ForkParsed(job.exec)
}
