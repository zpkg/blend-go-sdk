package main

import (
	"context"
	"flag"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
)

var exec = flag.String("exec", "", "The command to execute")
var name = flag.String("exec", stringutil.Letters.Random(8), "The command to execute")
var schedule = flag.String("schedule", "", "The job schedule")
var config = flag.String("config", "config.yml", "The job manager config")

func main() {
	flag.Parse()
}

// Job is the main job body.
type Job struct {
	StringSchedule *cron.StringSchedule
	Config         *jobkit.Config

	Name string
	Exec string
}

// Name returns the job name.
func (job Job) Name() string {
	return job.Config.Name
}

// Execute is the job body.
func (job Job) Execute(ctx context.Context) error {
	return sh.Exec(job.Exec)
}
