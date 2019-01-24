package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/blend/go-sdk/airbrake"
	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/diagnostics"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/stringutil"
)

var name = flag.String("name", stringutil.Letters.Random(8), "The name of the job")
var exec = flag.String("exec", "", "The command to execute")
var bind = flag.String("bind", "", "The address and port to bind the management server to (ex: 127.0.0.1:9000")
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

	// set up myriad of notification targets
	var emailClient email.Sender
	if !config.AWS.IsZero() {
		log.SyncInfof("email notifications: %s", logger.ColorBlue.Apply("enabled"))
		emailClient = ses.New(&config.AWS)
	}
	var slackClient slack.Sender
	if !config.Slack.IsZero() {
		log.SyncInfof("slack notifications: %s", logger.ColorBlue.Apply("enabled"))
		slackClient = slack.New(&config.Slack)
	}
	var statsClient stats.Collector
	if !config.Datadog.IsZero() {
		log.SyncInfof("datadog: %s", logger.ColorBlue.Apply("enabled"))
		statsClient, err = datadog.New(&config.Datadog)
		if err != nil {
			logger.FatalExit(err)
		}
	}
	var errorClient diagnostics.Notifier
	if !config.Airbrake.IsZero() {
		log.SyncInfof("airbrake: %s", logger.ColorBlue.Apply("enabled"))
		errorClient = airbrake.New(&config.Airbrake)
	}

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
		return sh.ForkContext(ctx, command[0], args(command...)...)
	}

	var notificationsConfig *jobkit.NotificationsConfig
	if cfg, ok := config.Notifications[*name]; ok {
		notificationsConfig = &cfg
	}

	jobs := cron.NewFromConfig(&config.Config).WithLogger(log)
	jobs.LoadJob(jobkit.NewJob(action).
		WithName(*name).
		WithSchedule(schedule).
		WithLogger(log).
		WithNotificationsConfig(notificationsConfig).
		WithEmailClient(emailClient).
		WithStatsClient(statsClient).
		WithSlackClient(slackClient).
		WithErrorClient(errorClient))

	ws := jobkit.NewManagementServer(jobs, &config)
	ws.WithLogger(log)
	if *bind != "" {
		ws.WithBindAddr(*bind)
	}

	go func() {
		if err := graceful.Shutdown(jobs); err != nil {
			logger.FatalExit(err)
		}
	}()
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
