package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type ActionsService interface {
	GetWorkflowJobByID(ctx context.Context, owner, repo string, jobID int64) (*WorkflowJob, *Response, error)
	GetWorkflowJobLogs(ctx context.Context, owner, repo string, jobID int64, maxRedirects int) (*url.URL, *Response, error)
	GetWorkflowRunByID(ctx context.Context, owner, repo string, runID int64) (*WorkflowRun, *Response, error)
	GetWorkflowRunAttempt(ctx context.Context, owner, repo string, runID int64, attemptNumber int, opts *WorkflowRunAttemptOptions) (*WorkflowRun, *Response, error)
	ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, opts *ListWorkflowJobsOptions) (*Jobs, *Response, error)
	ListWorkflowJobsAttempt(ctx context.Context, owner, repo string, runID, attemptNumber int64, opts *ListOptions) (*Jobs, *Response, error)
}

const maxPerPage = 100

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
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read error response body: %w", slogerr.With(err, "status_code", resp.StatusCode))
		}
		return nil, fmt.Errorf("download workflow job logs: %w", slogerr.With(errInvalidStatusCode, "status_code", resp.StatusCode, "response_body", string(b)))
	}
	return resp.Body, nil
}

var errInvalidStatusCode = errors.New("invalid status code")

func (c *Client) ListWorkflowJobs(ctx context.Context, owner, repo string, runID int64, attempt int) ([]*WorkflowJob, error) {
	list := c.getListJobsFunc(attempt)
	opts := &ListWorkflowJobsOptions{
		ListOptions: ListOptions{
			PerPage: maxPerPage,
		},
	}
	arr := []*WorkflowJob{}
	for range 10 { // max 1000 jobs
		jobs, resp, err := list(ctx, owner, repo, runID, opts.Page)
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

func (c *Client) getListJobsFunc(attempt int) func(ctx context.Context, owner, repo string, runID int64, page int) (*Jobs, *Response, error) {
	if attempt > 0 {
		return func(ctx context.Context, owner, repo string, runID int64, page int) (*Jobs, *Response, error) {
			return c.actions.ListWorkflowJobsAttempt(ctx, owner, repo, runID, int64(attempt), &ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
		}
	}
	return func(ctx context.Context, owner, repo string, runID int64, page int) (*Jobs, *Response, error) {
		return c.actions.ListWorkflowJobs(ctx, owner, repo, runID, &ListWorkflowJobsOptions{
			ListOptions: ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			},
		})
	}
}
