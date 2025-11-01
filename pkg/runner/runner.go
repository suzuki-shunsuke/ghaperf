package runner

import (
	"archive/zip"
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/view"
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
	ShowRun(run *collector.WorkflowRun, threshold time.Duration)
	ShowRuns(runs []*collector.WorkflowRun, threshold time.Duration)
}

type Collector interface {
	GetJobLog(ctx context.Context, input *collector.Input, jobID int64) ([]byte, error)
	GetJob(ctx context.Context, logger *slog.Logger, input *collector.Input, jobID int64) (*collector.Job, error)
	GetRun(ctx context.Context, logger *slog.Logger, input *collector.Input, runID int64, attempt int) (*collector.WorkflowRun, error)
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
