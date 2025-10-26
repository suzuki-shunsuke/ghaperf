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
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/log"
	"github.com/suzuki-shunsuke/ghaperf/pkg/runner"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type InputRun struct {
	LogLevel    string
	Repo        string
	RunID       int64
	JobID       int64
	StepName    string
	Threshold   string
	LogFile     string
	Args        []string
	EnableGHTKN bool
	Help        bool
	Version     bool
}

const (
	envLogLevel           = "GHAPERF_LOG_LEVEL"
	envGhaperfGitHubToken = "GHAPERF_GITHUB_TOKEN" //nolint:gosec
	envEnableGHTKN        = "GHAPERF_GHTKN"
	envGhaperfThreshold   = "GHAPERF_THRESHOLD"
	envGitHubRunID        = "GITHUB_RUN_ID"
	envGitHubRepository   = "GITHUB_REPOSITORY"
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

func (c *Controller) Run(ctx context.Context, logger *slog.Logger, logLevelVar *slog.LevelVar, input *InputRun, arg *Arg) error {
	if err := setLogLevel(logLevelVar, input.LogLevel, arg.Getenv); err != nil {
		return err
	}

	if err := setEnableGHTKN(input, arg.Getenv); err != nil {
		return err
	}

	threshold, err := getThreshold(input.Threshold, arg.Getenv)
	if err != nil {
		return err
	}

	if input.LogFile != "" {
		if err := runner.NewRunner(nil, arg.Stdout, arg.Fs).RunWithLogFile(logger, &collector.Input{
			Threshold: threshold,
			LogFile:   input.LogFile,
		}); err != nil {
			return fmt.Errorf("run with log file: %w", err)
		}
		return nil
	}

	job, err := getJobArg(input, arg)
	if err != nil {
		return err
	}

	gh, err := github.New(ctx, logger, &github.InputNew{
		GHTKNEnabled: input.EnableGHTKN,
		AccessToken:  getGitHubToken(arg.Getenv),
	})
	if err != nil {
		return fmt.Errorf("create GitHub client: %w", err)
	}
	if err := runner.NewRunner(gh, arg.Stdout, arg.Fs).Run(ctx, logger, &collector.Input{
		Threshold: threshold,
		Job:       job,
		CacheDir:  xdg.CacheDir(arg.Getenv, arg.Home),
	}); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

func getJobArg(input *InputRun, arg *Arg) (*collector.Job, error) {
	runID, err := getRunID(input.RunID, arg.Getenv)
	if err != nil {
		return nil, err
	}
	if runID == 0 && input.JobID == 0 {
		return nil, errors.New("one of --run-id, --job-id, and --log-file must be specified")
	}

	repoFullName := getRepoFullName(input.Repo, arg.Getenv)
	if repoFullName == "" {
		return nil, errors.New("without --log-file, repository must be specified")
	}

	repoOwner, repoName, err := validateRepo(repoFullName)
	if err != nil {
		return nil, err
	}
	return &collector.Job{
		JobID:     input.JobID,
		RunID:     runID,
		RepoOwner: repoOwner,
		RepoName:  repoName,
	}, nil
}

func getRepoFullName(s string, getEnv func(string) string) string {
	if s != "" {
		return s
	}
	return getEnv(envGitHubRepository)
}

const defaultThreshold = 30 * time.Second

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

func getRunID(runID int64, getEnv func(string) string) (int64, error) {
	if runID != 0 {
		return runID, nil
	}
	s := getEnv(envGitHubRunID)
	if s == "" {
		return 0, nil
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be int64: %w", envGitHubRunID, err)
	}
	return id, nil
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
