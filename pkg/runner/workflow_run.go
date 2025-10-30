package runner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
)

func (r *Runner) runWithRunID(ctx context.Context, logger *slog.Logger, input *collector.Input) error {
	run, jobs, err := r.collector.GetRun(ctx, logger, input, input.RunID, input.AttemptNumber)
	if err != nil {
		return fmt.Errorf("get jobs by run id: %w", err)
	}
	r.viewer.ShowJobs(run, jobs, input.Threshold)
	return nil
}

func (r *Runner) runs(ctx context.Context, logger *slog.Logger, input *collector.Input) error {
	runs, err := r.collector.ListRuns(ctx, logger, input, input.WorkflowNumber)
	if err != nil {
		return fmt.Errorf("list workflow runs: %w", err)
	}
	r.viewer.ShowRuns(runs, input.Threshold)
	return nil
}
