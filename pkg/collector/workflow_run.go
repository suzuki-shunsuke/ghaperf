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
)

const statusCompleted = "completed"

func (r *Collector) GetRun(ctx context.Context, logger *slog.Logger, input *Input) ([]*github.WorkflowJob, error) {
	jobIDsPath := xdg.RunJobIDsCache(input.CacheDir, input.Job.RepoOwner, input.Job.RepoName, input.Job.RunID)
	b, err := afero.ReadFile(r.fs, jobIDsPath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return nil, fmt.Errorf("read cached job IDs file: %w", err)
		}
		return r.getAndCacheRun(ctx, logger, input, jobIDsPath)
	}
	jobIDs := []int64{}
	// exist cache
	if err := json.Unmarshal(b, &jobIDs); err != nil {
		return nil, fmt.Errorf("unmarshal cached job ids: %w", err)
	}
	arr := make([]*github.WorkflowJob, len(jobIDs))
	for i, jobID := range jobIDs {
		jobCachePath := xdg.JobCache(input.CacheDir, input.Job.RepoOwner, input.Job.RepoName, jobID)
		job, err := r.getJob(ctx, logger, input, jobID, jobCachePath)
		if err != nil {
			return nil, fmt.Errorf("get a job: %w", err)
		}
		arr[i] = job
	}
	return arr, nil
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

func (r *Collector) getAndCacheRun(ctx context.Context, logger *slog.Logger, input *Input, jobIDsPath string) ([]*github.WorkflowJob, error) {
	run, err := r.gh.GetWorkflowRunByID(ctx, input.Job.RepoOwner, input.Job.RepoName, input.Job.RunID)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	// cache not found
	jobs, err := r.gh.ListWorkflowJobs(ctx, input.Job.RepoOwner, input.Job.RepoName, input.Job.RunID)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	if run.GetStatus() == statusCompleted {
		// cache workflow run and job ids
		if err := r.cacheJobIDs(jobs, jobIDsPath); err != nil {
			return nil, fmt.Errorf("cache job IDs: %w", err)
		}
		if err := r.cacheRun(run, xdg.RunCache(input.CacheDir, input.Job.RepoOwner, input.Job.RepoName, input.Job.RunID)); err != nil {
			return nil, fmt.Errorf("cache a workflow run: %w", err)
		}
	}
	// cache jobs
	if err := r.cacheJobs(logger, input, jobs); err != nil {
		return nil, fmt.Errorf("cache jobs: %w", err)
	}
	return jobs, nil
}

func (r *Collector) cacheJobs(logger *slog.Logger, input *Input, jobs []*github.WorkflowJob) error {
	for _, job := range jobs {
		if job.GetStatus() != statusCompleted {
			logger.Warn("job is not completed yet", "job_id", job.GetID(), "status", job.GetStatus())
			continue
		}
		jobCachePath := xdg.JobCache(input.CacheDir, input.Job.RepoOwner, input.Job.RepoName, job.GetID())
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
