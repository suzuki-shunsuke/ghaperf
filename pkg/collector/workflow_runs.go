package collector

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
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
	arr := make([]*WorkflowRun, len(runs))
	for i, run := range runs {
		logArgs := []any{"run_id", run.GetID()}
		cachePath := xdg.RunCache(input.CacheDir, input.RepoOwner, input.RepoName, run.GetID(), 0)
		if err := r.cacheRun(run, cachePath); err != nil {
			slogerr.WithError(logger, err).Error("cache workflow run", logArgs...)
		}
		jobs, err := r.getJobs(ctx, logger, input, run)
		if err != nil {
			slogerr.WithError(logger, err).Error("get jobs", logArgs...)
		}
		arr[i] = &WorkflowRun{
			Run:  run,
			Jobs: jobs,
		}
	}
	return arr, nil
}
