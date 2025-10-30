package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

const statusCompleted = "completed"

func (r *Collector) GetRun(ctx context.Context, logger *slog.Logger, input *Input) (*github.WorkflowRun, []*Job, error) {
	run, err := r.getRun(ctx, logger, input)
	if err != nil {
		return nil, nil, fmt.Errorf("get a workflow run: %w", err)
	}
	jobs, err := r.getJobs(ctx, logger, input, run)
	if err != nil {
		return nil, nil, fmt.Errorf("get jobs: %w", err)
	}
	return run, jobs, nil
}

func (r *Collector) cacheJobIDs(jobs []*github.WorkflowJob, cachePath string) error {
	jobIDs := make([]int64, len(jobs))
	for i, job := range jobs {
		jobIDs[i] = job.GetID()
	}
	b, err := json.Marshal(jobIDs)
	if err != nil {
		return fmt.Errorf("marshal job IDs: %w", err)
	}
	if err := r.fs.MkdirAll(filepath.Dir(cachePath), dirPermission); err != nil {
		return fmt.Errorf("make dirs for cached job IDs file: %w", err)
	}
	if err := afero.WriteFile(r.fs, cachePath, b, filePermission); err != nil {
		return fmt.Errorf("write cached job IDs file: %w", err)
	}
	return nil
}

func (r *Collector) cacheRun(workflowRun *github.WorkflowRun, cachePath string) error {
	b, err := json.Marshal(workflowRun)
	if err != nil {
		return fmt.Errorf("marshal workflow run: %w", err)
	}
	if err := r.fs.MkdirAll(filepath.Dir(cachePath), dirPermission); err != nil {
		return fmt.Errorf("make dirs for cached workflow run file: %w", err)
	}
	if err := afero.WriteFile(r.fs, cachePath, b, filePermission); err != nil {
		return fmt.Errorf("write cached workflow run file: %w", err)
	}
	return nil
}

func (r *Collector) getRun(ctx context.Context, logger *slog.Logger, input *Input) (*github.WorkflowRun, error) {
	cachePath := xdg.RunCache(input.CacheDir, input.RepoOwner, input.RepoName, input.RunID, input.AttemptNumber)
	runB, err := afero.ReadFile(r.fs, cachePath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			slogerr.WithError(logger, err).Error("read cached workflow run file")
		}
		return r.getAndCacheRun(ctx, logger, input, cachePath)
	}
	run := &github.WorkflowRun{}
	if err := json.Unmarshal(runB, run); err != nil {
		return nil, fmt.Errorf("unmarshal cached workflow run file: %w", err)
	}
	return run, nil
}

func (r *Collector) getAndCacheRun(ctx context.Context, logger *slog.Logger, input *Input, cachePath string) (*github.WorkflowRun, error) {
	run, err := r.gh.GetWorkflowRunByID(ctx, input.RepoOwner, input.RepoName, input.RunID, input.AttemptNumber)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	if run.GetStatus() == statusCompleted {
		// cache workflow run
		if err := r.cacheRun(run, cachePath); err != nil {
			slogerr.WithError(logger, err).Error("cache a workflow run")
		}
	}
	return run, nil
}
