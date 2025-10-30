package collector

import (
	"archive/zip"
	"context"
	"io"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
)

type Collector struct {
	fs afero.Fs
	gh GitHub
}

type GitHub interface {
	GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*github.WorkflowJob, error)
	GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64) (io.ReadCloser, error)
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64, attempt int) (*github.WorkflowRun, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, attempt int) ([]*github.WorkflowJob, error)
	ListWorkflowRuns(ctx context.Context, owner, repo string, fileName string, maxCount int, opts *github.ListWorkflowRunsOptions) ([]*github.WorkflowRun, error)
	GetWorkflowRunLogs(ctx context.Context, owner, repo string, runID int64, attempt int) ([]*zip.File, error)
}

func New(fs afero.Fs, gh GitHub) *Collector {
	return &Collector{
		fs: fs,
		gh: gh,
	}
}

type Input struct {
	Threshold               time.Duration
	LogFile                 string
	CacheDir                string
	RepoOwner               string
	RepoName                string
	RunID                   int64
	JobID                   int64
	AttemptNumber           int
	WorkflowNumber          int
	WorkflowName            string
	ListWorkflowRunsOptions *github.ListWorkflowRunsOptions
}

type Job struct {
	Job      *github.WorkflowJob
	Groups   []*parser.Group
	duration time.Duration
}

func (j *Job) Duration() time.Duration {
	if j == nil || j.Job == nil {
		return 0
	}
	if j.Job.GetConclusion() == "skipped" {
		return 0
	}
	if j.duration != 0 {
		return j.duration
	}
	completedAt := j.Job.GetCompletedAt().Time
	startedAt := j.Job.GetStartedAt().Time
	if completedAt.Before(startedAt) {
		return 0
	}
	j.duration = completedAt.Sub(startedAt)
	return j.duration
}
