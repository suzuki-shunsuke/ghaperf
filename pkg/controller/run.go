package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/config"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/log"
	"github.com/suzuki-shunsuke/ghaperf/pkg/runner"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type InputRun struct {
	LogLevel                string
	Repo                    string
	AttemptNumber           int
	RunID                   int64
	JobID                   int64
	Threshold               string
	LogFile                 string
	Args                    []string
	EnableGHTKN             bool
	Help                    bool
	Version                 bool
	ListWorkflowRunsOptions *github.ListWorkflowRunsOptions
	WorkflowNumber          int
	WorkflowName            string
	Config                  string
}

const (
	envLogLevel           = "GHAPERF_LOG_LEVEL"
	envGhaperfGitHubToken = "GHAPERF_GITHUB_TOKEN" //nolint:gosec
	envEnableGHTKN        = "GHAPERF_GHTKN"
	envGhaperfThreshold   = "GHAPERF_THRESHOLD"
	envGitHubToken        = "GITHUB_TOKEN" //nolint:gosec
)

type Arg struct {
	Getenv  func(string) string
	Stdout  io.Writer
	Stderr  io.Writer
	Fs      afero.Fs
	Args    []string
	Version string
	Home    string
}

func (c *Controller) Run(ctx context.Context, logger *slog.Logger, logLevelVar *slog.LevelVar, inputRun *InputRun, arg *Arg) error {
	if err := setLogLevel(logLevelVar, inputRun.LogLevel, arg.Getenv); err != nil {
		return err
	}
	input, err := c.getInput(inputRun, arg)
	if err != nil {
		return err
	}

	rArgs := &runner.Args{
		Stdout: arg.Stdout,
		Fs:     arg.Fs,
	}

	if inputRun.LogFile != "" {
		if err := runner.NewRunner(nil, rArgs).RunWithLogFile(input); err != nil {
			return fmt.Errorf("run with log file: %w", err)
		}
		return nil
	}

	gh, err := github.New(ctx, logger, &github.InputNew{
		GHTKNEnabled: inputRun.EnableGHTKN,
		AccessToken:  getGitHubToken(arg.Getenv),
	})
	if err != nil {
		return fmt.Errorf("create GitHub client: %w", err)
	}
	if err := runner.NewRunner(gh, rArgs).Run(ctx, logger, input); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

func (c *Controller) getInput(input *InputRun, arg *Arg) (*collector.Input, error) {
	if err := setEnableGHTKN(input, arg.Getenv); err != nil {
		return nil, err
	}

	threshold, err := getThreshold(input.Threshold, arg.Getenv)
	if err != nil {
		return nil, err
	}

	if input.LogFile != "" {
		return &collector.Input{
			Threshold: threshold,
			LogFile:   input.LogFile,
			Version:   arg.Version,
		}, nil
	}

	if input.RunID == 0 && input.JobID == 0 && input.WorkflowName == "" {
		return nil, errors.New("one of --run-id, --job-id, --log-file, and --workflow must be specified")
	}

	if input.Repo == "" {
		return nil, errors.New("without --log-file, repository must be specified")
	}

	repoOwner, repoName, err := validateRepo(input.Repo)
	if err != nil {
		return nil, err
	}

	cfg := &config.Config{}
	if err := readConfig(arg.Fs, input.Config, cfg); err != nil {
		return nil, err
	}

	return &collector.Input{
		Threshold:               threshold,
		CacheDir:                xdg.CacheDir(arg.Getenv, arg.Home),
		RepoOwner:               repoOwner,
		RepoName:                repoName,
		RunID:                   input.RunID,
		JobID:                   input.JobID,
		AttemptNumber:           input.AttemptNumber,
		WorkflowNumber:          input.WorkflowNumber,
		WorkflowName:            input.WorkflowName,
		ListWorkflowRunsOptions: input.ListWorkflowRunsOptions,
		Config:                  cfg,
		Version:                 arg.Version,
	}, nil
}

const defaultThreshold = 30 * time.Second

func readConfig(fs afero.Fs, path string, cfg *config.Config) error {
	if path == "" {
		return nil
	}
	if err := config.Read(fs, path, cfg); err != nil {
		return fmt.Errorf("read config file: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config file: %w", err)
	}
	return nil
}

func getThreshold(s string, getEnv func(string) string) (time.Duration, error) {
	threshold := getThresholdStr(s, getEnv)
	if threshold == "" {
		return defaultThreshold, nil
	}
	d, err := time.ParseDuration(threshold)
	if err != nil {
		return 0, fmt.Errorf("invalid threshold. See https://pkg.go.dev/time#ParseDuration: %w", err)
	}
	return d, nil
}

func getThresholdStr(s string, getEnv func(string) string) string {
	if s != "" {
		return s
	}
	return getEnv(envGhaperfThreshold)
}

func getGitHubToken(getEnv func(string) string) string {
	if token := getEnv(envGhaperfGitHubToken); token != "" {
		return token
	}
	return getEnv(envGitHubToken)
}

var errInvalidRepoArg = errors.New("invalid repository name format")

func validateRepo(repo string) (string, string, error) {
	repoOwner, repoName, ok := strings.Cut(repo, "/")
	if !ok {
		return "", "", slogerr.With(errInvalidRepoArg, "repository", repo) //nolint:wrapcheck
	}
	if strings.Contains(repoName, "/") {
		return "", "", slogerr.With(errInvalidRepoArg, "repository", repo) //nolint:wrapcheck
	}
	return repoOwner, repoName, nil
}

func setLogLevel(levelVar *slog.LevelVar, level string, getEnv func(string) string) error {
	if level == "" {
		level = getEnv(envLogLevel)
	}
	if level != "" {
		if err := log.SetLevel(levelVar, level); err != nil {
			return err //nolint:wrapcheck
		}
	}
	return nil
}

func setEnableGHTKN(flag *InputRun, getEnv func(string) string) error {
	if flag.EnableGHTKN {
		return nil
	}
	s := getEnv(envEnableGHTKN)
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("%s must be boolean: %w", envEnableGHTKN, err)
	}
	flag.EnableGHTKN = b
	return nil
}
