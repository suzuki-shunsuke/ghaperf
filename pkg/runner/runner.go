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
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64, attempt int) (*github.WorkflowRun, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, attempt int) ([]*github.WorkflowJob, error)
}

type Viewer interface {
	ShowJob(job *collector.Job, threshold time.Duration)
	ShowGroups(groups []*parser.Group, threshold time.Duration)
	ShowJobs(run *github.WorkflowRun, jobs []*collector.Job, threshold time.Duration)
}

type Collector interface {
	GetJobLog(ctx context.Context, input *collector.Input, jobID int64) ([]byte, error)
	GetJob(ctx context.Context, logger *slog.Logger, input *collector.Input, jobID int64) (*collector.Job, error)
	GetRun(ctx context.Context, logger *slog.Logger, input *collector.Input) (*github.WorkflowRun, []*collector.Job, error)
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
	if input.JobID != 0 {
		job, err := r.collector.GetJob(ctx, logger, input, input.JobID)
		if err != nil {
			return fmt.Errorf("run job ID %d: %w", input.JobID, err)
		}
		r.viewer.ShowJob(job, input.Threshold)
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
