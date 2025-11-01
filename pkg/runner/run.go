package runner

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/view"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (r *Runner) Run(ctx context.Context, logger *slog.Logger, input *collector.Input) error {
	headerArg := &view.HeaderArg{
		Version:                 input.Version,
		Repo:                    input.RepoOwner + "/" + input.RepoName,
		Now:                     time.Now(),
		Threshold:               input.Threshold,
		ListWorkflowRunsOptions: input.ListWorkflowRunsOptions,
		Count:                   input.WorkflowNumber,
		WorkflowName:            input.WorkflowName,
		Config:                  input.Config,
	}
	if input.JobID != 0 {
		job, err := r.collector.GetJob(ctx, logger, input, input.JobID)
		if err != nil {
			return fmt.Errorf("run job ID %d: %w", input.JobID, err)
		}
		r.viewer.ShowHeader(headerArg)
		r.viewer.ShowJob(job, input.Threshold)
		return nil
	}
	if input.RunID != 0 {
		return r.runWithRunID(ctx, logger, input, headerArg)
	}
	return r.runs(ctx, logger, input, headerArg)
}

func (r *Runner) RunWithLogFile(input *collector.Input) error {
	f, err := r.fs.Open(input.LogFile)
	if err != nil {
		return fmt.Errorf("open a log file: %w", slogerr.With(err, "log_file", input.LogFile))
	}
	defer f.Close()
	log, err := parser.Parse(f)
	if err != nil {
		return fmt.Errorf("parse a log file: %w", err)
	}
	r.viewer.ShowHeader(&view.HeaderArg{
		Version:                 input.Version,
		Now:                     time.Now(),
		Threshold:               input.Threshold,
		ListWorkflowRunsOptions: input.ListWorkflowRunsOptions,
	})
	r.viewer.ShowGroups(log.Groups, input.Threshold)
	return nil
}
