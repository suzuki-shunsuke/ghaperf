package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type ActionsService interface {
	GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*WorkflowJob, *Response, error)
	GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64, maxRedirects int) (*url.URL, *Response, error)
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64) (*WorkflowRun, *Response, error)
	GetWorkflowRunAttempt(ctx context.Context, owner, repo string, runID int64, attemptNumber int, opts *WorkflowRunAttemptOptions) (*WorkflowRun, *Response, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *ListWorkflowJobsOptions) (*Jobs, *Response, error)
	ListWorkflowJobsAttempt(ctx context.Context, owner, repo string, runID, attemptNumber int64, opts *ListOptions) (*Jobs, *Response, error)
	ListWorkflowRunsByFileName(ctx context.Context, owner, repo, workflowFileName string, opts *ListWorkflowRunsOptions) (*WorkflowRuns, *Response, error)
	GetWorkflowRunLogs(ctx context.Context, owner, repo string, runID int64, maxRedirects int) (*url.URL, *Response, error)
	GetWorkflowRunAttemptLogs(ctx context.Context, owner, repo string, runID int64, attemptNumber, maxRedirects int) (*url.URL, *Response, error)
}

const (
	maxPerPage   = 100
	maxRedirects = 5
)

var errInvalidStatusCode = errors.New("invalid status code")

func (c *Client) GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*WorkflowJob, error) {
	job, _, err := c.actions.GetWorkflowJobByID(ctx, owner, repo, jobID)
	if err != nil {
		return nil, fmt.Errorf("get workflow job by ID: %w", err)
	}
	return job, nil
}

func (c *Client) GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64, attempt int) (*WorkflowRun, error) {
	if attempt > 0 {
		run, _, err := c.actions.GetWorkflowRunAttempt(ctx, owner, repo, runID, attempt, nil)
		if err != nil {
			return nil, fmt.Errorf("get workflow run attempt by ID: %w", err)
		}
		return run, nil
	}
	run, _, err := c.actions.GetWorkflowRunByID(ctx, owner, repo, runID)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	return run, nil
}
