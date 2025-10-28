package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/pflag"
	"github.com/suzuki-shunsuke/ghaperf/pkg/controller"
)

const help = `ghaperf - Analyze the performance of GitHub Actions using GitHub API and raw job logs
https://github.com/suzuki-shunsuke/ghaperf

USAGE:
   ghaperf --help [-h] # Show this help
   ghaperf --version [-v] # Show version
   ghaperf [OPTIONS]

OPTIONS:
   --log-level <debug|info|warn|error>            Set log level (or set GHIR_LOG_LEVEL)
   --ghtkn                                        Enable the integration with ghtkn (or set GHAPERF_GHTKN)
   --repo <owner>/<repo>                          Specify the repository
   --run-id <run id>                              Specify the run ID
   --job-id <job id>                              Specify the job ID
   --step-name <regular expression of step name>  Specify the step name
   --threshold <time duration> 				      Specify the threshold duration (e.g., 30s, 1m)
   --log-file <file path>						  Specify the job log file path
   --help, -h                                     Show help
   --version, -v                                  Show version

VERSION:
   %s
`

func parseFlags(f *controller.InputRun) {
	pflag.StringVar(&f.LogLevel, "log-level", "", "log level (debug, info, warn, error)")
	pflag.StringVar(&f.Repo, "repo", "", "repository (owner/repo)")
	pflag.Int64Var(&f.RunID, "run-id", 0, "run ID")
	pflag.Int64Var(&f.JobID, "job-id", 0, "job ID")
	pflag.StringVar(&f.StepName, "step-name", "", "step name")
	pflag.StringVar(&f.Threshold, "threshold", "", "threshold")
	pflag.StringVar(&f.LogFile, "log-file", "", "log file")
	pflag.BoolVar(&f.EnableGHTKN, "ghtkn", false, "Enable the integration with ghtkn")
	pflag.BoolVarP(&f.Help, "help", "h", false, "Show help")
	pflag.BoolVarP(&f.Version, "version", "v", false, "Show version")
	pflag.Parse()
	f.Args = pflag.Args()
}

func Run(ctx context.Context, logger *slog.Logger, logLevel *slog.LevelVar, arg *controller.Arg) error {
	inputRun := &controller.InputRun{}
	parseFlags(inputRun)

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
