package collector

import (
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
}

func New(fs afero.Fs, gh GitHub) *Collector {
	return &Collector{
		fs: fs,
		gh: gh,
	}
}

type Input struct {
	Threshold     time.Duration
	LogFile       string
	Data          string
	CacheDir      string
	RepoOwner     string
	RepoName      string
	RunID         int64
	JobID         int64
	AttemptNumber int
}

type Job struct {
	Job    *github.WorkflowJob
	Groups []*parser.Group
}
