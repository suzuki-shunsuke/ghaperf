package xdg

import (
	"path/filepath"
	"strconv"
)

const (
	envXDGCacheHome = "XDG_CACHE_HOME"
)

func CacheDir(getEnv func(string) string, home string) string {
	return filepath.Join(cacheHome(getEnv, home), "ghaperf")
}

func cacheHome(getEnv func(string) string, home string) string {
	if s := getEnv(envXDGCacheHome); s != "" {
		return s
	}
	return filepath.Join(home, ".cache")
}

func JobCache(cacheDir, repoOwner, repoName string, jobID int64) string {
	return filepath.Join(cacheDir, "jobs", repoOwner, repoName, strconv.FormatInt(jobID, 10), "job.json")
}

func JobLogCache(jobCachePath string) string {
	return filepath.Join(filepath.Dir(jobCachePath), "log.txt")
}

func RunCache(cacheDir, repoOwner, repoName string, runID int64, attempt int) string {
	file := "run.json"
	dir := filepath.Join(cacheDir, "runs", repoOwner, repoName, strconv.FormatInt(runID, 10))
	if attempt == 0 {
		return filepath.Join(dir, file)
	}
	return filepath.Join(dir, strconv.Itoa(attempt), file)
}

func RunJobIDsCache(cacheDir, repoOwner, repoName string, runID int64, attempt int) string {
	file := "job_ids.json"
	dir := filepath.Join(cacheDir, "runs", repoOwner, repoName, strconv.FormatInt(runID, 10))
	if attempt == 0 {
		return filepath.Join(dir, file)
	}
	return filepath.Join(dir, strconv.Itoa(attempt), file)
}
