package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/uuid"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
)

var name = flag.String("name", stringutil.Letters.Random(8), "The name of the job")
var exec = flag.String("exec", "", "The command to execute")
var bind = flag.String("bind", "", "The address and port to bind the management server to")
var schedule = flag.String("schedule", "*/1 * * * * * *", "The job schedule as a cron string (i.e. 7 space delimited components)")
var configPath = flag.String("config", "config.yml", "The job config path")
var timeout = flag.Duration("timeout", 0, "The timeout")

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

	var command []string
	if *exec != "" {
		command = strings.Split(*exec, " ")
	} else {
		command, err = sh.ParseFlagsTrailer(os.Args...)
		if err != nil {
			logger.FatalExit(err)
		}
	}

	if len(command) == 0 {
		logger.FatalExit(fmt.Errorf("must supply a command to run with `--exec=...` or `-- command`)"))
	}

	jm := cron.New().WithLogger(log)
	jm.LoadJob(&ExecJob{
		schedule: schedule,
		name:     *name,
		exec:     command[0],
		args:     args(command...),
		timeout:  *timeout,
	})

	go func() {
		if err := graceful.Shutdown(jm); err != nil {
			logger.FatalExit(err)
		}
	}()

	ws := jobkit.NewManagementServer(jm, &config)
	ws.WithLogger(log)
	if *bind != "" {
		ws.WithBindAddr(*bind)
	}
	if err := graceful.Shutdown(ws); err != nil {
		logger.FatalExit(err)
	}
}

func args(all ...string) []string {
	if len(all) < 2 {
		return nil
	}
	return all[1:]
}

// NewExecJob creates a new exec job.
func NewExecJob(exec string, args ...string) *ExecJob {
	return &ExecJob{
		name: uuid.V4().String(),
		exec: exec,
		args: args,
	}
}

var (
	_ cron.Job                    = (*ExecJob)(nil)
	_ cron.OnStartReceiver        = (*ExecJob)(nil)
	_ cron.OnCompleteReceiver     = (*ExecJob)(nil)
	_ cron.OnFailureReceiver      = (*ExecJob)(nil)
	_ cron.OnCancellationReceiver = (*ExecJob)(nil)
	_ cron.OnBrokenReceiver       = (*ExecJob)(nil)
	_ cron.OnFixedReceiver        = (*ExecJob)(nil)
)

// ExecJob is the main job body.
type ExecJob struct {
	schedule cron.Schedule
	config   *jobkit.Config
	name     string
	exec     string
	args     []string
	timeout  time.Duration
}

// Name returns the job name.
func (job ExecJob) Name() string {
	return job.name
}

// WithName sets the name.
func (job *ExecJob) WithName(name string) *ExecJob {
	job.name = name
	return job
}

// Schedule returns the job schedule.
func (job ExecJob) Schedule() cron.Schedule {
	return job.schedule
}

// WithSchedule sets the schedule.
func (job *ExecJob) WithSchedule(schedule cron.Schedule) *ExecJob {
	job.schedule = schedule
	return job
}

// Timeout returns the timeout.
func (job ExecJob) Timeout() time.Duration {
	return job.timeout
}

// WithTimeout sets the job timeout.
func (job *ExecJob) WithTimeout(d time.Duration) *ExecJob {
	job.timeout = d
	return job
}

// OnStart is a lifecycle event handler.
func (job ExecJob) OnStart(ctx context.Context) {
	//
}

// OnComplete is a lifecycle event handler.
func (job ExecJob) OnComplete(ctx context.Context) {
	//
}

// OnFailure is a lifecycle event handler.
func (job ExecJob) OnFailure(ctx context.Context) {
	//
}

// OnBroken is a lifecycle event handler.
func (job ExecJob) OnBroken(ctx context.Context) {
	//
}

// OnFixed is a lifecycle event handler.
func (job ExecJob) OnFixed(ctx context.Context) {
	//
}

// OnCancellation is a lifecycle event handler.
func (job ExecJob) OnCancellation(ctx context.Context) {
	//
}

// Exec returns the job command.
func (job ExecJob) Exec() string {
	return job.exec
}

// Args returns the job command args.
func (job ExecJob) Args() []string {
	return job.args
}

// Execute is the job body.
func (job ExecJob) Execute(ctx context.Context) error {
	return sh.ForkContext(ctx, job.exec, job.args...)
}
