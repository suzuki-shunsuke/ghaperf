package runner

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/view"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (r *Runner) runWithRunID(ctx context.Context, logger *slog.Logger, input *collector.Input, headerArg *view.HeaderArg) error {
	run, err := r.collector.GetRun(ctx, logger, input, input.RunID, input.AttemptNumber)
	if err != nil {
		if !errors.Is(err, github.ErrLogHasGone) {
			return fmt.Errorf("get run by run id: %w", err)
		}
		slogerr.WithError(logger, err).Warn("get run by run id")
	}
	r.viewer.ShowHeader(headerArg)
	r.viewer.ShowJobs(run, input.Threshold)
	return nil
}

func (r *Runner) runs(ctx context.Context, logger *slog.Logger, input *collector.Input, headerArg *view.HeaderArg) error {
	runs, err := r.collector.ListRuns(ctx, logger, input, input.WorkflowNumber)
	if err != nil {
		return fmt.Errorf("list workflow runs: %w", err)
	}
	r.viewer.ShowHeader(headerArg)
	r.viewer.ShowRuns(runs, input.Threshold)
	return nil
}
