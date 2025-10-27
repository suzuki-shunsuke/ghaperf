package collector

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
)

func (c *Collector) GetJobLog(ctx context.Context, input *Input, jobID int64) ([]byte, error) {
	cachePath := xdg.JobLogCache(xdg.JobCache(input.CacheDir, input.RepoOwner, input.RepoName, jobID))
	b, err := afero.ReadFile(c.fs, cachePath)
	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return nil, fmt.Errorf("read cached job log file: %w", err)
		}
		// cache not found
		return c.getAndCacheLog(ctx, input, jobID, cachePath)
	}
	// exist cache
	return b, nil
}

const filePermission = 0o644

func (c *Collector) getAndCacheLog(ctx context.Context, input *Input, jobID int64, cachePath string) ([]byte, error) {
	logReader, err := c.gh.GetWorkflowJobLogs(ctx, input.RepoOwner, input.RepoName, jobID)
	if err != nil {
		return nil, fmt.Errorf("get workflow job logs: %w", err)
	}
	defer logReader.Close()
	b, err := io.ReadAll(logReader)
	if err != nil {
		return nil, fmt.Errorf("read workflow job logs: %w", err)
	}
	// cache the job info
	if err := afero.WriteFile(c.fs, cachePath, b, filePermission); err != nil {
		return nil, fmt.Errorf("write cached job log file: %w", err)
	}
	return b, nil
}
