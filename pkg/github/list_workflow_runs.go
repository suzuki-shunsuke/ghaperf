package github

import (
	"context"
	"fmt"
)

func (c *Client) ListWorkflowRuns(ctx context.Context, owner, repo string, fileName string, maxCount int, opts *ListWorkflowRunsOptions) ([]*WorkflowRun, error) {
	o := &ListWorkflowRunsOptions{
		Actor:               opts.Actor,
		Branch:              opts.Branch,
		Event:               opts.Event,
		Status:              opts.Status,
		Created:             opts.Created,
		HeadSHA:             opts.HeadSHA,
		ExcludePullRequests: opts.ExcludePullRequests,
		CheckSuiteID:        opts.CheckSuiteID,
		ListOptions: ListOptions{
			PerPage: maxPerPage,
		},
	}
	if maxCount < maxPerPage {
		o.PerPage = maxCount
	}
	arr := []*WorkflowRun{}
	for range 10 { // max 1000 jobs
		runs, resp, err := c.actions.ListWorkflowRunsByFileName(ctx, owner, repo, fileName, o)
		if err != nil {
			return nil, fmt.Errorf("list workflow runs: %w", err)
		}
		arr = append(arr, runs.WorkflowRuns...)
		if resp.NextPage == 0 {
			return arr, nil
		}
		if len(arr) >= maxCount {
			return arr[:maxCount], nil
		}
		o.Page = resp.NextPage
	}
	return arr, nil
}
