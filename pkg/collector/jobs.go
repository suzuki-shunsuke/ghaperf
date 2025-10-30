package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (r *Collector) getJobs(ctx context.Context, logger *slog.Logger, input *Input, run *github.WorkflowRun) ([]*Job, error) {
	cachePath := xdg.RunJobIDsCache(input.CacheDir, input.RepoOwner, input.RepoName, run.GetID(), input.AttemptNumber)
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
	arr := make([]*Job, len(jobIDs))
	for i, jobID := range jobIDs {
		job, err := r.GetJob(ctx, logger, input, jobID)
		if err != nil {
			return nil, fmt.Errorf("get a job: %w", err)
		}
		arr[i] = job
	}
	return arr, nil
}

func (r *Collector) getAndCacheJobs(ctx context.Context, logger *slog.Logger, input *Input, jobIDsPath string, run *github.WorkflowRun) ([]*Job, error) {
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
	arr := make([]*Job, 0, len(jobs))
	for _, job := range jobs {
		logArgs := []any{"job_id", job.GetID(), "job_name", job.GetName(), "job_status", job.GetStatus()}
		if job.GetStatus() != statusCompleted {
			logger.Warn("exclude a not completed job", logArgs...)
			continue
		}
		jobLog, err := r.GetJobLog(ctx, input, job.GetID())
		if err != nil {
			slogerr.WithError(logger, err).Error("get a job log", logArgs...)
			arr = append(arr, &Job{
				Job: job,
			})
			continue
		}
		groups, err := parser.Parse(logger, bytes.NewBuffer(jobLog))
		if err != nil {
			slogerr.WithError(logger, err).Error("parse a job log", logArgs...)
		}
		arr = append(arr, &Job{
			Job:    job,
			Groups: groups,
		})
	}
	return arr, nil
}

func (r *Collector) cacheJobs(logger *slog.Logger, input *Input, jobs []*github.WorkflowJob) error {
	for _, job := range jobs {
		if job.GetStatus() != statusCompleted {
			logger.Warn("job is not completed yet", "job_id", job.GetID(), "status", job.GetStatus())
			continue
		}
		jobCachePath := xdg.JobCache(input.CacheDir, input.RepoOwner, input.RepoName, job.GetID())
		if f, err := afero.Exists(r.fs, jobCachePath); err == nil && f {
			continue
		}
		// cache the job info
		if err := r.cacheJob(jobCachePath, job); err != nil {
			return fmt.Errorf("cache job info: %w", err)
		}
	}
	return nil
}
