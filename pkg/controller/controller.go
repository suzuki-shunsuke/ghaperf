package controller

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
)

type Controller struct {
	input *InputNew
}

func New(input *InputNew) *Controller {
	return &Controller{
		input: input,
	}
}

type GitHub interface {
	GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*github.WorkflowJob, error)
	GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64) (io.ReadCloser, error)
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64) (*github.WorkflowRun, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64) ([]*github.WorkflowJob, error)
}

type InputNew struct{}

func NewInput() *InputNew {
	return &InputNew{}
}
