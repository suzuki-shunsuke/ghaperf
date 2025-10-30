package collector

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type WorkflowRun struct {
	Run  *github.WorkflowRun
	Jobs []*Job
}

func (r *Collector) ListRuns(ctx context.Context, logger *slog.Logger, input *Input, maxCount int) ([]*WorkflowRun, error) {
	runs, err := r.gh.ListWorkflowRuns(ctx, input.RepoOwner, input.RepoName, input.WorkflowName, maxCount, input.ListWorkflowRunsOptions)
	if err != nil {
		return nil, fmt.Errorf("list workflow runs: %w", err)
	}
	arr := make([]*WorkflowRun, 0, len(runs))
	for _, run := range runs {
		logArgs := []any{"run_id", run.GetID(), "run_attempt", run.GetRunAttempt()}
		jobs, err := r.getJobsAndLogs(ctx, logger, input, run)
		if err != nil {
			slogerr.WithError(logger, err).Error("get jobs and logs", logArgs...)
			continue
		}
		arr = append(arr, &WorkflowRun{
			Run:  run,
			Jobs: jobs,
		})
	}
	return arr, nil
}
