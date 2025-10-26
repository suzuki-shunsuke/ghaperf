package controller

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
)

func (r *Runner) getJobLog(ctx context.Context, input *Input, jobID int64, jobCachePath string) ([]byte, error) {
	cachePath := xdg.JobLogCache(jobCachePath)
	b, err := afero.ReadFile(r.fs, cachePath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return nil, fmt.Errorf("read cached job log file: %w", err)
		}
		// cache not found
		return r.getAndCacheLog(ctx, input, jobID, cachePath)
	}
	// exist cache
	return b, nil
}

const filePermission = 0o644

func (r *Runner) getAndCacheLog(ctx context.Context, input *Input, jobID int64, cachePath string) ([]byte, error) {
	logReader, err := r.gh.GetWorkflowJobLogs(ctx, input.Job.RepoOwner, input.Job.RepoName, jobID)
	if err != nil {
		return nil, fmt.Errorf("get workflow job logs: %w", err)
	}
	defer logReader.Close()
	b, err := io.ReadAll(logReader)
	if err != nil {
		return nil, fmt.Errorf("read workflow job logs: %w", err)
	}
	// cache the job info
	if err := afero.WriteFile(r.fs, cachePath, b, filePermission); err != nil {
		return nil, fmt.Errorf("write cached job log file: %w", err)
	}
	return b, nil
}
