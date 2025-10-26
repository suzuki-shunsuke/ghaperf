package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/view"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type Runner struct {
	gh     GitHub
	stdout io.Writer
	fs     afero.Fs
	viewer Viewer
}

type Viewer interface {
	ShowJob(groups []*parser.Group, threshold time.Duration, job *github.WorkflowJob)
	ShowGroups(groups []*parser.Group, threshold time.Duration)
	ShowJobs(jobs []*github.WorkflowJob, threshold time.Duration)
}

func NewRunner(gh GitHub, stdout io.Writer, fs afero.Fs) *Runner {
	return &Runner{
		gh:     gh,
		stdout: stdout,
		fs:     fs,
		viewer: view.New(stdout),
	}
}

type Input struct {
	Threshold time.Duration
	LogFile   string
	Data      string
	Job       *Job
	CacheDir  string
}

func (r *Runner) Run(ctx context.Context, logger *slog.Logger, input *Input) error {
	if input.Job.JobID != 0 {
		job, groups, err := r.runJob(ctx, logger, input, input.Job.JobID)
		if err != nil {
			return fmt.Errorf("run job ID %d: %w", input.Job.JobID, err)
		}
		r.viewer.ShowJob(groups, input.Threshold, job)
		return nil
	}
	return r.RunWithRunID(ctx, logger, input)
}

func (r *Runner) RunWithLogFile(logger *slog.Logger, input *Input) error {
	f, err := r.fs.Open(input.LogFile)
	if err != nil {
		return fmt.Errorf("open a log file: %w", slogerr.With(err, "log_file", input.LogFile))
	}
	defer f.Close()
	groups, err := parser.Parse(logger, f)
	if err != nil {
		return fmt.Errorf("parse a log file: %w", err)
	}
	r.viewer.ShowGroups(groups, input.Threshold)
	return nil
}

func (r *Runner) runJob(ctx context.Context, logger *slog.Logger, input *Input, jobID int64) (*github.WorkflowJob, []*parser.Group, error) {
	jobCachePath := xdg.JobCache(input.CacheDir, input.Job.RepoOwner, input.Job.RepoName, jobID)
	job, err := r.getJob(ctx, logger, input, jobID, jobCachePath)
	if err != nil {
		return nil, nil, fmt.Errorf("get a job: %w", err)
	}
	jobLog, err := r.getJobLog(ctx, input, jobID, jobCachePath)
	if err != nil {
		return nil, nil, fmt.Errorf("get a job log: %w", err)
	}
	groups, err := parser.Parse(logger, bytes.NewBuffer(jobLog))
	if err != nil {
		return nil, nil, fmt.Errorf("parse log: %w", err)
	}
	return job, groups, nil
}
