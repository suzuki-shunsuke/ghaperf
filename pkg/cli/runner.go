package cli

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/spf13/pflag"
	"github.com/suzuki-shunsuke/ghaperf/pkg/controller"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
)

const help = `ghaperf - Analyze the performance of GitHub Actions using GitHub API and raw job logs
https://github.com/suzuki-shunsuke/ghaperf

USAGE:
   ghaperf --help [-h] # Show this help
   ghaperf --version [-v] # Show version
   ghaperf [OPTIONS]

OPTIONS:
   --log-level <debug|info|warn|error>    Log level (or set GHIR_LOG_LEVEL)
   --ghtkn                                Enable the integration with ghtkn (or set GHAPERF_GHTKN)
   --repo <owner>/<repo>                  The repository
   --run-id <run id>                      The run ID
   --job-id <job id>                      The job ID
   --attempt-number <attempt number>      The workflow run's attempt number
   --threshold <time duration>            The threshold duration (e.g., 30s, 1m)
   --log-file <file path>                 Log file path
   --count <the number of workflow runs>  The number of workflow runs to analyze (default: 100)
   --workflow <workflow name>             The workflow name
   --workflow-actor <actor>               The workflow run actor
   --workflow-branch <branch>             The workflow run branch
   --workflow-event <event>               The workflow run event
   --workflow-created <date range>        The workflow run created date range
   --workflow-status <status>             The workflow run status
   --config <path>                        The config file path
   --init                                 Initialize the config file
   --help, -h                             Show help
   --version, -v                          Show version

VERSION:
   %s
`

func parseFlags(f *controller.InputRun) {
	if f.ListWorkflowRunsOptions == nil {
		f.ListWorkflowRunsOptions = &github.ListWorkflowRunsOptions{}
	}
	pflag.StringVar(&f.LogLevel, "log-level", "", "log level (debug, info, warn, error)")
	pflag.StringVar(&f.Repo, "repo", "", "repository (owner/repo)")
	pflag.Int64Var(&f.RunID, "run-id", 0, "run ID")
	pflag.IntVar(&f.AttemptNumber, "attempt-number", 0, "workflow run's attempt number")
	pflag.Int64Var(&f.JobID, "job-id", 0, "job ID")
	pflag.StringVar(&f.Threshold, "threshold", "", "threshold")
	pflag.StringVar(&f.LogFile, "log-file", "", "log file")
	pflag.BoolVar(&f.EnableGHTKN, "ghtkn", false, "Enable the integration with ghtkn")
	pflag.BoolVarP(&f.Help, "help", "h", false, "Show help")
	pflag.BoolVar(&f.Init, "init", false, "Initialize the config file")
	pflag.BoolVarP(&f.Version, "version", "v", false, "Show version")
	pflag.IntVar(&f.WorkflowNumber, "count", 100, "the number of workflow runs") //nolint:mnd
	pflag.StringVar(&f.WorkflowName, "workflow", "", "the workflow name")
	pflag.StringVar(&f.ListWorkflowRunsOptions.Actor, "workflow-actor", "", "the workflow run actor")
	pflag.StringVar(&f.ListWorkflowRunsOptions.Branch, "workflow-branch", "", "the workflow run branch")
	pflag.StringVar(&f.ListWorkflowRunsOptions.Event, "workflow-event", "", "the workflow run event")
	pflag.StringVar(&f.ListWorkflowRunsOptions.Created, "workflow-created", "", "the workflow run created date range")
	pflag.StringVar(&f.ListWorkflowRunsOptions.Status, "workflow-status", "", "the workflow run status")
	pflag.StringVar(&f.Config, "config", "", "the config file path")

	pflag.Parse()
	f.Args = pflag.Args()
}

var digitPrefix = regexp.MustCompile("^[0-9]")

func Run(ctx context.Context, logger *slog.Logger, logLevel *slog.LevelVar, arg *controller.Arg) error {
	inputRun := &controller.InputRun{}
	parseFlags(inputRun)

	if digitPrefix.MatchString(arg.Version) {
		arg.Version = "v" + arg.Version
	}

	if inputRun.Help {
		fmt.Fprintf(arg.Stdout, help, arg.Version)
		return nil
	}
	if inputRun.Version {
		fmt.Fprintf(arg.Stdout, "%s\n", arg.Version)
		return nil
	}

	ctrl := controller.New(&controller.InputNew{})
	if err := ctrl.Run(ctx, logger, logLevel, inputRun, arg); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}
