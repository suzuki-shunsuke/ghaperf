package github

import (
	"context"
	"fmt"
)

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
