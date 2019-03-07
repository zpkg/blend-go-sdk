package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/blend/go-sdk/airbrake"
	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/diagnostics"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/stringutil"
)

// the following flags apply to any invocation
var bind = flag.String("bind", "", "The address and port to bind the management server to (ex: 127.0.0.1:9000")
var configPath = flag.String("config", "config.yml", "The job config path")
var disableServer = flag.Bool("disable-server", false, "Disables the management server (will make --bind irrelevant)")

// the following flags create a default job
var defaultJobName = flag.String("name", "", "The name of the job")
var defaultJobExec = flag.String("exec", "", "The command to execute")
var defaultJobSchedule = flag.String("schedule", "", "The job schedule as a cron string (i.e. 7 space delimited components)")
var defaultJobDiscardOutput = flag.Bool("discard-output", false, "Discard job output")
var defaultJobTimeout = flag.Duration("timeout", 0, "The timeout")

type config struct {
	jobkit.Config `json:",inline" yaml:",inline"`
	DisableServer *bool `json:"disableServer" yaml:"disableServer"`

	Jobs []jobConfig `json:"jobs" yaml:"jobs"`
}

func (c *config) Resolve() error {
	if err := configutil.SetString(&c.Web.BindAddr, configutil.String(*bind), configutil.Env("BIND_ADDR"), configutil.String(c.Web.BindAddr)); err != nil {
		return err
	}
	if err := configutil.SetBool(&c.DisableServer, configutil.Bool(disableServer), configutil.Bool(c.DisableServer), configutil.Bool(ref.Bool(false))); err != nil {
		return err
	}
	return nil
}

type jobConfig struct {
	Exec          string   `json:"exec" yaml:"exec"`
	Command       []string `json:"command" yaml:"command"`
	DiscardOutput *bool    `json:"discardOutput" yaml:"discardOutput"`

	jobkit.JobConfig `json:",inline" yaml:",inline"`
}

func (jc *jobConfig) Resolve() error {
	return configutil.AnyError(
		configutil.SetString(&jc.Name, configutil.String(*defaultJobName), configutil.String(env.Env().ServiceName()), configutil.String(jc.Name), configutil.String(stringutil.Letters.Random(8))),
		configutil.SetString(&jc.Exec, configutil.String(*defaultJobExec), configutil.String(jc.Exec)),
		configutil.SetStrings(&jc.Command, configutil.StringsFunc(argsTrailer), configutil.Strings(jc.Command)),
		configutil.SetBool(&jc.DiscardOutput, configutil.Bool(defaultJobDiscardOutput), configutil.Bool(jc.DiscardOutput), configutil.Bool(ref.Bool(false))),
		configutil.SetString(&jc.Schedule, configutil.String(*defaultJobSchedule), configutil.String(jc.Schedule)),
		configutil.SetDuration(&jc.Timeout, configutil.Duration(*defaultJobTimeout), configutil.Duration(jc.Timeout)),
	)
}

func argsTrailer() ([]string, error) {
	command, _ := sh.ArgsTrailer(os.Args...)
	if len(command) == 0 {
		return nil, nil
	}
	return command, nil
}

func main() {
	flag.Parse()

	var err error
	var cfg config
	if err := configutil.Read(&cfg, *configPath); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}

	log := logger.NewFromConfig(&cfg.Logger)
	log.WithEnabled(cron.FlagStarted, cron.FlagComplete, cron.FlagFixed, cron.FlagBroken, cron.FlagFailed, cron.FlagCancelled)

	defaultJobCfg, err := createDefaultJobConfig()
	if err != nil {
		log.SyncFatalExit(err)
	}
	if defaultJobCfg != nil {
		cfg.Jobs = append(cfg.Jobs, *defaultJobCfg)
	}

	if len(cfg.Jobs) == 0 {
		logger.FatalExit(fmt.Errorf("must supply a command to run with `--exec=...` or `-- command`), or provide a jobs config file"))
	}

	// set up myriad of notification targets
	var emailClient email.Sender
	if !cfg.AWS.IsZero() {
		emailClient = ses.New(aws.MustNewSession(&cfg.AWS))
	}
	var slackClient slack.Sender
	if !cfg.Slack.IsZero() {
		slackClient = slack.New(&cfg.Slack)
	}
	var statsClient stats.Collector
	if !cfg.Datadog.IsZero() {
		statsClient, err = datadog.NewCollector(&cfg.Datadog)
		if err != nil {
			log.SyncFatalExit(err)
		}
	}
	var errorClient diagnostics.Notifier
	if !cfg.Airbrake.IsZero() {
		errorClient = airbrake.MustNew(&cfg.Airbrake)
	}

	jobs := cron.NewFromConfig(&cfg.Config.Config).WithLogger(log)

	for _, jobCfg := range cfg.Jobs {
		job, err := createJob(&jobCfg)
		if err != nil {
			log.SyncFatalExit(err)
		}
		job.WithLogger(log).WithEmailClient(emailClient).WithSlackClient(slackClient).WithStatsClient(statsClient).WithErrorClient(errorClient)
		log.SyncInfof("loading job `%s` with schedule `%s`", jobCfg.NameOrDefault(), jobCfg.ScheduleOrDefault())
		jobs.LoadJob(job)
	}

	if !*disableServer {
		ws := jobkit.NewManagementServer(jobs, &cfg.Config).WithLogger(log)
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

func createDefaultJobConfig() (*jobConfig, error) {
	cfg := new(jobConfig)
	if err := cfg.Resolve(); err != nil {
		return nil, err
	}
	if cfg.Exec == "" && len(cfg.Command) == 0 {
		return nil, nil
	}
	return cfg, nil
}

func createJob(cfg *jobConfig) (*jobkit.Job, error) {
	if cfg.Exec == "" && len(cfg.Command) == 0 {
		return nil, exception.New("job exec and command unset").WithMessagef("job: %s", cfg.NameOrDefault())
	}
	var command []string
	if cfg.Exec != "" {
		command = stringutil.SplitSpaceQuoted(cfg.Exec)
	} else {
		command = cfg.Command
	}
	action := func(ctx context.Context) error {
		if cfg.DiscardOutput == nil || (cfg.DiscardOutput != nil && !*cfg.DiscardOutput) {
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

	job, err := jobkit.NewJob(&cfg.JobConfig, action)
	if err != nil {
		return nil, err
	}
	if job.Description() == "" {
		job.WithDescription(strings.Join(command, " "))
	}
	return job, nil
}

func args(all ...string) []string {
	if len(all) < 2 {
		return nil
	}
	return all[1:]
}
