package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/stringutil"
)

var name = flag.String("name", "", "The name of the job")
var exec = flag.String("exec", "", "The command to execute")
var bind = flag.String("bind", "", "The address and port to bind the management server to (ex: 127.0.0.1:9000")
var schedule = flag.String("schedule", "", "The job schedule as a cron string (i.e. 7 space delimited components)")
var disableServer = flag.Bool("disable-server", false, "Disables the management server (will make --bind irrelevant)")
var discardOutput = flag.Bool("discard-output", false, "Discard job output")
var configPath = flag.String("config", "config.yml", "The job config path")
var timeout = flag.Duration("timeout", 0, "The timeout")

type jobConfig struct {
	jobkit.Config    `json:",inline" yaml:",inline"`
	jobkit.JobConfig `json:",inline" yaml:",inline"`
}

func (jc *jobConfig) Resolve() error {
	return configutil.AnyError(
		configutil.SetString(&jc.Name, configutil.String(*name), configutil.String(env.Env().ServiceName()), configutil.String(jc.Name), configutil.String(stringutil.Letters.Random(8))),
		configutil.SetString(&jc.Schedule, configutil.String(*schedule), configutil.String(jc.Schedule)),
		configutil.SetDuration(&jc.Timeout, configutil.Duration(*timeout), configutil.Duration(jc.Timeout)),
		configutil.SetString(&jc.Web.BindAddr, configutil.String(*bind), configutil.String(jc.Web.BindAddr)),
	)
}

func main() {
	flag.Parse()

	var err error
	var config jobConfig
	if err := configutil.Read(&config, *configPath); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}

	log := logger.NewFromConfig(&config.Logger)
	log.WithEnabled(cron.FlagStarted, cron.FlagComplete, cron.FlagFixed, cron.FlagBroken, cron.FlagFailed, cron.FlagCancelled)

	log.SyncInfof("starting job `%s`", config.NameOrDefault())
	log.SyncInfof("using schedule `%s`", config.ScheduleOrDefault())

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

	action := func(ctx context.Context) error {
		if !*discardOutput {
			if jis := jobkit.GetJobInvocationState(ctx); jis != nil {
				cmd, err := sh.CmdContext(ctx, command[0], args(command...)...)
				if err != nil {
					return err
				}
				cmd.Stdout = io.MultiWriter(jis.Output, os.Stdout)
				cmd.Stderr = io.MultiWriter(jis.ErrorOutput, os.Stderr)
				return exception.New(cmd.Run())
			}
		}
		return sh.ForkContext(ctx, command[0], args(command...)...)
	}

	job, err := jobkit.NewJob(&config.JobConfig, &config.Config, action)
	if err != nil {
		logger.FatalExit(err)
	}
	job.WithLogger(log)

	jobs := cron.NewFromConfig(&config.Config.Config).WithLogger(log)
	jobs.LoadJob(job)

	if !*disableServer {
		ws := jobkit.NewManagementServer(jobs, &config.Config).WithLogger(log)
		go func() {
			if err := graceful.Shutdown(ws); err != nil {
				logger.FatalExit(err)
			}
		}()
	}

	if err := graceful.Shutdown(graceful.New(jobs.Start, jobs.Stop)); err != nil {
		logger.FatalExit(err)
	}
}

func args(all ...string) []string {
	if len(all) < 2 {
		return nil
	}
	return all[1:]
}
