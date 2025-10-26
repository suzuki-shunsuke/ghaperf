package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
)

func (r *Runner) getJob(ctx context.Context, logger *slog.Logger, input *Input, jobID int64, jobCachePath string) (*github.WorkflowJob, error) {
	job := &github.WorkflowJob{}
	b, err := afero.ReadFile(r.fs, jobCachePath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return nil, fmt.Errorf("read cached job file: %w", err)
		}
		// cache not found
		job, err := r.gh.GetWorkflowJobByID(ctx, input.Job.RepoOwner, input.Job.RepoName, jobID)
		if err != nil {
			return nil, fmt.Errorf("get workflow job by ID: %w", err)
		}
		if job.GetStatus() != "completed" {
			logger.Warn("job is not completed yet", "job_id", jobID, "status", job.GetStatus())
			return job, nil
		}
		// cache the job info
		if err := r.cacheJob(jobCachePath, job); err != nil {
			return nil, fmt.Errorf("cache job info: %w", err)
		}
		return job, nil
	}
	// exist cache
	if err := json.Unmarshal(b, job); err != nil {
		return nil, fmt.Errorf("unmarshal cached job file: %w", err)
	}
	return job, nil
}

const dirPermission = 0o755

func (r *Runner) cacheJob(jobCachePath string, job *github.WorkflowJob) error {
	if err := r.fs.MkdirAll(filepath.Dir(jobCachePath), dirPermission); err != nil {
		return fmt.Errorf("create job cache dir: %w", err)
	}
	b, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job file: %w", err)
	}
	if err := afero.WriteFile(r.fs, jobCachePath, b, filePermission); err != nil {
		return fmt.Errorf("write cached job file: %w", err)
	}
	return nil
}
