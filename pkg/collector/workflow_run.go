package collector

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"path/filepath"
	"slices"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

const statusCompleted = "completed"

func (r *Collector) GetRun(ctx context.Context, logger *slog.Logger, input *Input, runID int64, attempt int) (*github.WorkflowRun, []*Job, error) {
	run, err := r.getRun(ctx, logger, input, runID, attempt)
	if err != nil {
		return nil, nil, fmt.Errorf("get a workflow run: %w", err)
	}
	jobs, err := r.getJobsAndLogs(ctx, logger, input, run)
	if err != nil {
		return nil, nil, fmt.Errorf("get jobs: %w", err)
	}
	return run, jobs, nil
}

func (r *Collector) getJobsAndLogs(ctx context.Context, logger *slog.Logger, input *Input, run *github.WorkflowRun) ([]*Job, error) {
	jobs, err := r.getJobs(ctx, logger, input, run)
	if err != nil {
		return nil, fmt.Errorf("get jobs: %w", err)
	}
	jobM := make(map[string]*Job, len(jobs))
	for _, job := range jobs {
		name := input.Config.NormalizeJobName(job.GetName())
		if !input.Config.Include(job.GetName()) {
			continue
		}
		jobM[job.GetName()] = &Job{
			Job:            job,
			NormalizedName: name,
		}
	}
	logCacheDir := xdg.RunLogCache(input.CacheDir, input.RepoOwner, input.RepoName, run.GetID(), run.GetRunAttempt())
	logCacheFile := xdg.RunLogCacheFile(input.CacheDir, input.RepoOwner, input.RepoName, run.GetID(), run.GetRunAttempt())
	if f, err := afero.Exists(r.fs, logCacheFile); err == nil && f {
		// exist cache
		if err := r.readCachedLog(logger, logCacheDir, jobM); err != nil {
			return nil, fmt.Errorf("read cached logs: %w", err)
		}
		return slices.Collect(maps.Values(jobM)), nil
	}

	if err := r.fs.MkdirAll(logCacheDir, dirPermission); err != nil {
		return nil, fmt.Errorf("make dirs for cached workflow run log dir: %w", err)
	}
	files, err := r.gh.GetWorkflowRunLogs(ctx, input.RepoOwner, input.RepoName, run.GetID(), run.GetRunAttempt())
	if err != nil {
		return nil, fmt.Errorf("get workflow run logs: %w", err)
	}
	r.cacheAndParseLogs(logger, files, logCacheDir, logCacheFile, jobM)
	return slices.Collect(maps.Values(jobM)), nil
}

func (r *Collector) cacheAndParseLogs(logger *slog.Logger, files []*zip.File, logCacheDir, logCacheFile string, jobM map[string]*Job) {
	allCached := true
	for _, file := range files {
		cached, err := r.cacheAndParseLog(filepath.Join(logCacheDir, file.Name), file, jobM) //nolint:gosec
		if err != nil {
			slogerr.WithError(logger, err).Error("parse a cached log file", "file_name", file.Name)
		}
		if !cached {
			allCached = false
		}
	}
	if allCached {
		if err := afero.WriteFile(r.fs, logCacheFile, []byte{}, filePermission); err != nil {
			slogerr.WithError(logger, err).Error("write cached workflow run log file")
		}
	}
}

func (r *Collector) cacheAndParseLog(cachePath string, file *zip.File, jobM map[string]*Job) (bool, error) {
	if err := r.cacheLog(cachePath, file); err != nil {
		return false, err
	}
	log, err := r.readLog(cachePath)
	if err != nil {
		return true, err
	}
	job, ok := jobM[log.JobName]
	if !ok {
		return true, nil
	}
	job.Groups = log.Groups
	return true, nil
}

func (r *Collector) readCachedLog(logger *slog.Logger, logCacheDir string, jobM map[string]*Job) error {
	infos, err := afero.ReadDir(r.fs, logCacheDir)
	if err != nil {
		return fmt.Errorf("read cached workflow run log dir: %w", err)
	}
	for _, info := range infos {
		log, err := r.readLog(filepath.Join(logCacheDir, info.Name()))
		if err != nil {
			slogerr.WithError(logger, err).Error("parse a cached log file", "file_name", info.Name())
			continue
		}
		job, ok := jobM[log.JobName]
		if !ok {
			continue
		}
		job.Groups = log.Groups
	}
	return nil
}

func (r *Collector) cacheLog(cachePath string, file *zip.File) error {
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("open a log file from workflow run logs: %w", err)
	}
	defer f.Close()
	cacheFile, err := r.fs.Create(cachePath)
	if err != nil {
		return fmt.Errorf("create a cached log file: %w", err)
	}
	defer cacheFile.Close()
	// cache log
	if err := copySafe(cacheFile, f); err != nil {
		return fmt.Errorf("cache a log file: %w", err)
	}
	return nil
}

func (r *Collector) readLog(cachePath string) (*parser.Log, error) {
	f, err := r.fs.Open(cachePath)
	if err != nil {
		return nil, fmt.Errorf("open a cached log file: %w", err)
	}
	defer f.Close()
	log, err := parser.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parse a cached log file: %w", err)
	}
	return log, nil
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

func (r *Collector) getRun(ctx context.Context, logger *slog.Logger, input *Input, runID int64, attempt int) (*github.WorkflowRun, error) {
	var cachePath string
	if attempt > 0 {
		cachePath = xdg.RunCache(input.CacheDir, input.RepoOwner, input.RepoName, runID, attempt)
		runB, err := afero.ReadFile(r.fs, cachePath)
		if err != nil {
			if !errors.Is(err, afero.ErrFileNotFound) {
				slogerr.WithError(logger, err).Error("read cached workflow run file")
			}
			return r.getAndCacheRun(ctx, logger, input, runID, attempt)
		}
		run := &github.WorkflowRun{}
		if err := json.Unmarshal(runB, run); err != nil {
			return nil, fmt.Errorf("unmarshal cached workflow run file: %w", err)
		}
		return run, nil
	}
	return r.getAndCacheRun(ctx, logger, input, runID, attempt)
}

func (r *Collector) getAndCacheRun(ctx context.Context, logger *slog.Logger, input *Input, runID int64, attempt int) (*github.WorkflowRun, error) {
	run, err := r.gh.GetWorkflowRunByID(ctx, input.RepoOwner, input.RepoName, runID, attempt)
	if err != nil {
		return nil, fmt.Errorf("get workflow run by ID: %w", err)
	}
	if run.GetStatus() == statusCompleted {
		// cache workflow run
		if err := r.cacheRun(run, xdg.RunCache(input.CacheDir, input.RepoOwner, input.RepoName, input.RunID, run.GetRunAttempt())); err != nil {
			slogerr.WithError(logger, err).Error("cache a workflow run")
		}
	}
	return run, nil
}
