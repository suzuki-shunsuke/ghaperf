package runner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
)

func (r *Runner) runWithRunID(ctx context.Context, logger *slog.Logger, input *collector.Input) error {
	run, jobs, err := r.collector.GetRun(ctx, logger, input)
	if err != nil {
		return fmt.Errorf("get jobs by run id: %w", err)
	}
	r.viewer.ShowJobs(run, jobs, input.Threshold)
	return nil
}
