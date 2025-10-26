package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (r *Runner) RunWithRunID(ctx context.Context, logger *slog.Logger, input *Input) error {
	jobs, err := r.getRun(ctx, logger, input)
	if err != nil {
		return fmt.Errorf("get jobs by run id: %w", err)
	}
	arr := make([]*JobWithSteps, 0, len(jobs))
	for _, job := range jobs {
		job, slowSteps, err := r.runJob(ctx, logger, input, job.GetID())
		if err != nil {
			return fmt.Errorf("run job ID: %w", slogerr.With(err, "job_id", job.GetID()))
		}
		d := jobDuration(job)
		if d < input.Threshold {
			continue
		}
		arr = append(arr, &JobWithSteps{
			Job:       job,
			SlowSteps: slowSteps,
			Duration:  d,
		})
	}
	if len(arr) == 0 {
		fmt.Fprintln(r.stdout, "No slow jobs found")
		return nil
	}
	fmt.Fprintln(r.stdout, "## Slow jobs")
	for _, job := range arr {
		fmt.Fprintf(r.stdout, "Job Name: %s\n", job.Job.GetName())
		fmt.Fprintf(r.stdout, "Job ID: %d\n", job.Job.GetID())
		fmt.Fprintf(r.stdout, "Job Status: %s\n\n", job.Job.GetStatus())
		fmt.Fprintf(r.stdout, "Job Duration: %s\n\n", job.Duration)
		fmt.Fprintln(r.stdout, "### Slow steps")
		if len(job.SlowSteps) == 0 {
			fmt.Fprintf(r.stdout, "The job %s has no slow steps\n", job.Job.GetName())
			continue
		}
		for i, step := range job.SlowSteps {
			fmt.Fprintf(r.stdout, "%d. %s: %s\n", i+1, step.Duration, step.Name)
			for j, group := range step.Groups {
				fmt.Fprintf(r.stdout, "   %d %s: %s\n", j+1, group.Duration, group.Name)
			}
		}
	}
	return nil
}

const statusCompleted = "completed"

func (r *Runner) getRun(ctx context.Context, logger *slog.Logger, input *Input) ([]*github.WorkflowJob, error) {
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

func jobDuration(job *github.WorkflowJob) time.Duration {
	return job.GetCompletedAt().Sub(job.GetStartedAt().Time)
}

type JobWithSteps struct {
	Job       *github.WorkflowJob
	SlowSteps []*Step
	Duration  time.Duration
}

func (r *Runner) cacheJobIDs(jobs []*github.WorkflowJob, cachePath string) error {
	jobIDs := make([]int64, len(jobs))
	for i, job := range jobs {
		jobIDs[i] = job.GetID()
	}
	b, err := json.Marshal(jobIDs)
	if err != nil {
		return fmt.Errorf("marshal job IDs: %w", err)
	}
	if err := afero.WriteFile(r.fs, cachePath, b, filePermission); err != nil {
		return fmt.Errorf("write cached job IDs file: %w", err)
	}
	return nil
}

func (r *Runner) getAndCacheRun(ctx context.Context, logger *slog.Logger, input *Input, jobIDsPath string) ([]*github.WorkflowJob, error) {
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
		// cache jobs ids of workflow run
		if err := r.cacheJobIDs(jobs, jobIDsPath); err != nil {
			return nil, fmt.Errorf("cache job IDs: %w", err)
		}
	}
	// cache jobs
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
			return nil, fmt.Errorf("cache job info: %w", err)
		}
	}
	return jobs, nil
}
