package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ActionsService interface {
	GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*WorkflowJob, *Response, error)
	GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64, maxRedirects int) (*url.URL, *Response, error)
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64) (*WorkflowRun, *Response, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *ListWorkflowJobsOptions) (*Jobs, *Response, error)
}

func (c *Client) GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*WorkflowJob, error) {
	job, _, err := c.actions.GetWorkflowJobByID(ctx, owner, repo, jobID)
	if err != nil {
		return nil, fmt.Errorf("get workflow job by ID: %w", err)
	}
	return job, nil
}

func (c *Client) GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64) (*WorkflowRun, error) {
	run, _, err := c.actions.GetWorkflowRunByID(ctx, owner, repo, runID)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	return run, nil
}

const maxRedirects = 5

// GetWorkflowJobLogs downloads the logs for a workflow job.
// It returns an io.ReadCloser which must be closed by the caller.
func (c *Client) GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64) (io.ReadCloser, error) {
	link, _, err := c.actions.GetWorkflowJobLogs(ctx, owner, repo, jobID, maxRedirects)
	if err != nil {
		return nil, fmt.Errorf("get workflow job logs redirect URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create a http request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download workflow job logs: %w", err)
	}
	return resp.Body, nil
}

func (c *Client) ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64) ([]*WorkflowJob, error) {
	opts := &ListWorkflowJobsOptions{
		ListOptions: ListOptions{
			PerPage: 100, //nolint:mnd
		},
	}
	arr := []*WorkflowJob{}
	for range 10 { // max 1000 jobs
		jobs, resp, err := c.actions.ListWorkflowJobs(ctx, owner, repo, runID, opts)
		if err != nil {
			return nil, fmt.Errorf("list workflow jobs: %w", err)
		}
		arr = append(arr, jobs.Jobs...)
		if resp.NextPage == 0 {
			return arr, nil
		}
		opts.Page = resp.NextPage
	}
	return arr, nil
}
