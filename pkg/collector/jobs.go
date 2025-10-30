package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (r *Collector) getJobs(ctx context.Context, logger *slog.Logger, input *Input, run *github.WorkflowRun) ([]*github.WorkflowJob, error) {
	cachePath := xdg.RunJobIDsCache(input.CacheDir, input.RepoOwner, input.RepoName, run.GetID(), run.GetRunAttempt())
	b, err := afero.ReadFile(r.fs, cachePath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return nil, fmt.Errorf("read cached job IDs file: %w", err)
		}
		return r.getAndCacheJobs(ctx, logger, input, cachePath, run)
	}
	jobIDs := []int64{}
	// exist cache
	if err := json.Unmarshal(b, &jobIDs); err != nil {
		return nil, fmt.Errorf("unmarshal cached job ids: %w", err)
	}
	arr := make([]*github.WorkflowJob, len(jobIDs))
	for i, jobID := range jobIDs {
		job, err := r.getJob(ctx, logger, input, jobID)
		if err != nil {
			return nil, fmt.Errorf("get a job: %w", err)
		}
		arr[i] = job
	}
	return arr, nil
}

func (r *Collector) getAndCacheJobs(ctx context.Context, logger *slog.Logger, input *Input, jobIDsPath string, run *github.WorkflowRun) ([]*github.WorkflowJob, error) {
	// cache not found
	jobs, err := r.gh.ListWorkflowJobs(ctx, input.RepoOwner, input.RepoName, run.GetID(), input.AttemptNumber)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	if run.GetStatus() == statusCompleted {
		// cache workflow run and job ids
		if err := r.cacheJobIDs(jobs, jobIDsPath); err != nil {
			slogerr.WithError(logger, err).Error("cache job IDs")
		}
	}
	// cache jobs
	if err := r.cacheJobs(logger, input, jobs); err != nil {
		slogerr.WithError(logger, err).Error("cache jobs")
	}
	return jobs, nil
}

func (r *Collector) cacheJobs(logger *slog.Logger, input *Input, jobs []*github.WorkflowJob) error {
	for _, job := range jobs {
		logArgs := []any{"job_id", job.GetID(), "job_name", job.GetName(), "job_status", job.GetStatus(), "job_conclusion", job.GetConclusion()}
		if job.GetStatus() != statusCompleted {
			logger.Warn("job is not completed yet", logArgs...)
			continue
		}
		jobCachePath := xdg.JobCache(input.CacheDir, input.RepoOwner, input.RepoName, job.GetID())
		if f, err := afero.Exists(r.fs, jobCachePath); err == nil && f {
			continue
		}
		// cache the job info
		if err := r.cacheJob(jobCachePath, job); err != nil {
			return fmt.Errorf("cache job info: %w", slogerr.With(err, logArgs...))
		}
	}
	return nil
}
