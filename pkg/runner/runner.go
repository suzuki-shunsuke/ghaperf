package runner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/view"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type Runner struct {
	gh        GitHub
	stdout    io.Writer
	fs        afero.Fs
	viewer    Viewer
	collector Collector
}

type GitHub interface {
	GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*github.WorkflowJob, error)
	GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64) (io.ReadCloser, error)
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64) (*github.WorkflowRun, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64) ([]*github.WorkflowJob, error)
}

type Viewer interface {
	ShowJob(groups []*parser.Group, threshold time.Duration, job *github.WorkflowJob)
	ShowGroups(groups []*parser.Group, threshold time.Duration)
	ShowJobs(jobs []*github.WorkflowJob, threshold time.Duration)
}

type Collector interface {
	GetJobLog(ctx context.Context, input *collector.Input, jobID int64, jobCachePath string) ([]byte, error)
	GetJob(ctx context.Context, logger *slog.Logger, input *collector.Input, jobID int64) (*github.WorkflowJob, []*parser.Group, error)
	GetRun(ctx context.Context, logger *slog.Logger, input *collector.Input) ([]*github.WorkflowJob, error)
}

func NewRunner(gh GitHub, stdout io.Writer, fs afero.Fs) *Runner {
	return &Runner{
		gh:        gh,
		stdout:    stdout,
		fs:        fs,
		viewer:    view.New(stdout),
		collector: collector.New(fs, gh),
	}
}

func (r *Runner) Run(ctx context.Context, logger *slog.Logger, input *collector.Input) error {
	if input.Job.JobID != 0 {
		job, groups, err := r.collector.GetJob(ctx, logger, input, input.Job.JobID)
		if err != nil {
			return fmt.Errorf("run job ID %d: %w", input.Job.JobID, err)
		}
		r.viewer.ShowJob(groups, input.Threshold, job)
		return nil
	}
	return r.runWithRunID(ctx, logger, input)
}

func (r *Runner) RunWithLogFile(logger *slog.Logger, input *collector.Input) error {
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
