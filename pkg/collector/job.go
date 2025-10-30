package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Collector) GetJob(ctx context.Context, logger *slog.Logger, input *Input, jobID int64) (*Job, error) {
	job, err := c.getJob(ctx, logger, input, jobID)
	if err != nil {
		return nil, fmt.Errorf("get a job: %w", err)
	}
	logArgs := []any{"job_name", job.GetName(), "job_status", job.GetStatus(), "job_conclusion", job.GetConclusion()}
	if job.GetStatus() != statusCompleted {
		logger.Warn("job is not completed yet", logArgs...)
		return &Job{
			Job: job,
		}, nil
	}
	if job.GetConclusion() == "skipped" {
		return &Job{
			Job: job,
		}, nil
	}
	jobLog, err := c.GetJobLog(ctx, input, jobID)
	if err != nil {
		slogerr.WithError(logger, err).Error("get a job log", logArgs...)
		return &Job{
			Job: job,
		}, nil
	}
	groups, err := parser.Parse(logger, bytes.NewBuffer(jobLog))
	if err != nil {
		slogerr.WithError(logger, err).Error("parse a job log", logArgs...)
		return &Job{
			Job:    job,
			Groups: groups,
		}, nil
	}
	return &Job{
		Job:    job,
		Groups: groups,
	}, nil
}

func (c *Collector) getJob(ctx context.Context, logger *slog.Logger, input *Input, jobID int64) (*github.WorkflowJob, error) {
	jobCachePath := xdg.JobCache(input.CacheDir, input.RepoOwner, input.RepoName, jobID)
	job := &github.WorkflowJob{}
	b, err := afero.ReadFile(c.fs, jobCachePath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return nil, fmt.Errorf("read cached job file: %w", err)
		}
		// cache not found
		return c.getAndCacheJob(ctx, logger, input, jobID, jobCachePath)
	}
	// exist cache
	if err := json.Unmarshal(b, job); err != nil {
		return nil, fmt.Errorf("unmarshal cached job file: %w", err)
	}
	return job, nil
}

func (c *Collector) getAndCacheJob(ctx context.Context, logger *slog.Logger, input *Input, jobID int64, jobCachePath string) (*github.WorkflowJob, error) {
	job, err := c.gh.GetWorkflowJobByID(ctx, input.RepoOwner, input.RepoName, jobID)
	if err != nil {
		return nil, fmt.Errorf("get workflow job by ID: %w", err)
	}
	if job.GetStatus() != "completed" {
		logger.Warn("job is not completed yet", "job_id", jobID, "status", job.GetStatus())
		return job, nil
	}
	// cache the job info
	if err := c.cacheJob(jobCachePath, job); err != nil {
		return nil, fmt.Errorf("cache job info: %w", err)
	}
	return job, nil
}

const dirPermission = 0o755

func (c *Collector) cacheJob(jobCachePath string, job *github.WorkflowJob) error {
	if err := c.fs.MkdirAll(filepath.Dir(jobCachePath), dirPermission); err != nil {
		return fmt.Errorf("create job cache dir: %w", err)
	}
	b, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job file: %w", err)
	}
	if err := afero.WriteFile(c.fs, jobCachePath, b, filePermission); err != nil {
		return fmt.Errorf("write cached job file: %w", err)
	}
	return nil
}
