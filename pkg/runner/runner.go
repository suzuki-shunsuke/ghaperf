package runner

import (
	"archive/zip"
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
	ListWorkflowRuns(ctx context.Context, owner, repo string, fileName string, maxCount int, opts *github.ListWorkflowRunsOptions) ([]*github.WorkflowRun, error)
	GetWorkflowRunLogs(ctx context.Context, owner, repo string, runID int64, attempt int) ([]*zip.File, error)
}

type Viewer interface {
	ShowHeader(arg *view.HeaderArg)
	ShowJob(job *collector.Job, threshold time.Duration)
	ShowGroups(groups []*parser.Group, threshold time.Duration)
	ShowJobs(run *github.WorkflowRun, jobs []*collector.Job, threshold time.Duration)
	ShowRuns(runs []*collector.WorkflowRun, threshold time.Duration)
}

type Collector interface {
	GetJobLog(ctx context.Context, input *collector.Input, jobID int64) ([]byte, error)
	GetJob(ctx context.Context, logger *slog.Logger, input *collector.Input, jobID int64) (*collector.Job, error)
	GetRun(ctx context.Context, logger *slog.Logger, input *collector.Input, runID int64, attempt int) (*github.WorkflowRun, []*collector.Job, error)
	ListRuns(ctx context.Context, logger *slog.Logger, input *collector.Input, maxCount int) ([]*collector.WorkflowRun, error)
}

type Args struct {
	Stdout io.Writer
	Fs     afero.Fs
}

func NewRunner(gh GitHub, args *Args) *Runner {
	return &Runner{
		gh:        gh,
		stdout:    args.Stdout,
		fs:        args.Fs,
		viewer:    view.New(args.Stdout),
		collector: collector.New(args.Fs, gh),
	}
}

func (r *Runner) Run(ctx context.Context, logger *slog.Logger, input *collector.Input) error {
	headerArg := &view.HeaderArg{
		Version:                 input.Version,
		Now:                     time.Now(),
		Threshold:               input.Threshold,
		ListWorkflowRunsOptions: input.ListWorkflowRunsOptions,
		Count:                   input.WorkflowNumber,
		WorkflowName:            input.WorkflowName,
	}
	if input.JobID != 0 {
		job, err := r.collector.GetJob(ctx, logger, input, input.JobID)
		if err != nil {
			return fmt.Errorf("run job ID %d: %w", input.JobID, err)
		}
		r.viewer.ShowHeader(headerArg)
		r.viewer.ShowJob(job, input.Threshold)
		return nil
	}
	if input.RunID != 0 {
		return r.runWithRunID(ctx, logger, input, headerArg)
	}
	return r.runs(ctx, logger, input, headerArg)
}

func (r *Runner) RunWithLogFile(input *collector.Input) error {
	f, err := r.fs.Open(input.LogFile)
	if err != nil {
		return fmt.Errorf("open a log file: %w", slogerr.With(err, "log_file", input.LogFile))
	}
	defer f.Close()
	log, err := parser.Parse(f)
	if err != nil {
		return fmt.Errorf("parse a log file: %w", err)
	}
	r.viewer.ShowHeader(&view.HeaderArg{
		Version:                 input.Version,
		Now:                     time.Now(),
		Threshold:               input.Threshold,
		ListWorkflowRunsOptions: input.ListWorkflowRunsOptions,
	})
	r.viewer.ShowGroups(log.Groups, input.Threshold)
	return nil
}
