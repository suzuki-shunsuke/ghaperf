package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"time"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/xdg"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type Runner struct {
	gh     GitHub
	stdout io.Writer
	fs     afero.Fs
}

func NewRunner(gh GitHub, stdout io.Writer, fs afero.Fs) *Runner {
	return &Runner{
		gh:     gh,
		stdout: stdout,
		fs:     fs,
	}
}

type Input struct {
	Threshold time.Duration
	LogFile   string
	Data      string
	Job       *Job
	CacheDir  string
}

type Step struct {
	Name      string        `json:"name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Groups    []*Group      `json:"groups"`
}

func (s *Step) Contain(group *Group, threshold time.Duration) {
	if s.EndTime.Add(500 * time.Millisecond).Before(group.StartTime) { //nolint:mnd
		// The step ends before the group starts
		// Go to the next step
		return
	}
	if group.Duration < threshold {
		// The group is not slow
		// Go to the next group
		return
	}
	if s.StartTime.Add(-500*time.Millisecond).Before(group.StartTime) && s.EndTime.Add(500*time.Millisecond).After(group.EndTime) { //nolint:mnd
		// The step overlaps with the group
		s.Groups = append(s.Groups, group)
		return
	}
}

func (r *Runner) Run(ctx context.Context, logger *slog.Logger, input *Input) error {
	if input.Job.JobID != 0 {
		job, slowSteps, err := r.runJob(ctx, logger, input, input.Job.JobID)
		if err != nil {
			return fmt.Errorf("run job ID %d: %w", input.Job.JobID, err)
		}
		fmt.Fprintln(r.stdout, "## Slow steps")
		fmt.Fprintf(r.stdout, "Job Name: %s\n", job.GetName())
		fmt.Fprintf(r.stdout, "Job ID: %d\n", input.Job.JobID)
		fmt.Fprintf(r.stdout, "Job Status: %s\n\n", job.GetStatus())
		fmt.Fprintf(r.stdout, "Job Duration: %s\n\n", jobDuration(job))
		if len(slowSteps) == 0 {
			fmt.Fprintf(r.stdout, "The job %s has no slow steps\n", job.GetName())
			return nil
		}
		for i, step := range slowSteps {
			fmt.Fprintf(r.stdout, "%d. %s: %s\n", i+1, step.Duration, step.Name)
			for j, group := range step.Groups {
				fmt.Fprintf(r.stdout, "   %d %s: %s\n", j+1, group.Duration, group.Name)
			}
		}
		return nil
	}
	return r.RunWithRunID(ctx, logger, input)
}

func (r *Runner) RunWithLogFile(logger *slog.Logger, input *Input) error {
	f, err := r.fs.Open(input.LogFile)
	if err != nil {
		return fmt.Errorf("open a log file: %w", slogerr.With(err, "log_file", input.LogFile))
	}
	defer f.Close()
	groups, err := parseLog(logger, f)
	if err != nil {
		return err
	}
	slowGroups := make([]*Group, 0, len(groups))
	for _, group := range groups {
		if group.Duration < input.Threshold {
			continue
		}
		slowGroups = append(slowGroups, group)
	}
	if len(slowGroups) == 0 {
		fmt.Fprintln(r.stdout, "No slow log group is found")
		return nil
	}
	sort.Slice(slowGroups, func(i, j int) bool {
		return slowGroups[i].Duration > slowGroups[j].Duration
	})
	fmt.Fprintln(r.stdout, "## Slog log groups")
	for i, group := range slowGroups {
		fmt.Fprintf(r.stdout, "%d. %s: %s\n", i+1, group.Duration, group.Name)
	}
	return nil
}

func (r *Runner) runJob(ctx context.Context, logger *slog.Logger, input *Input, jobID int64) (*github.WorkflowJob, []*Step, error) {
	jobCachePath := xdg.JobCache(input.CacheDir, input.Job.RepoOwner, input.Job.RepoName, jobID)
	job, err := r.getJob(ctx, logger, input, jobID, jobCachePath)
	if err != nil {
		return nil, nil, fmt.Errorf("get a job: %w", err)
	}
	jobLog, err := r.getJobLog(ctx, input, jobID, jobCachePath)
	if err != nil {
		return nil, nil, fmt.Errorf("get a job log: %w", err)
	}
	groups, err := parseLog(logger, bytes.NewBuffer(jobLog))
	if err != nil {
		return nil, nil, fmt.Errorf("parse log: %w", err)
	}
	slowSteps := getSlowSteps(job.Steps, input.Threshold)
	if len(slowSteps) == 0 {
		return nil, nil, nil
	}
	slowGroups := getSlowGroups(groups, input.Threshold)
	for _, step := range slowSteps {
		for _, g := range slowGroups {
			step.Contain(g, input.Threshold)
		}
		sort.Slice(step.Groups, func(i, j int) bool {
			return step.Groups[i].Duration > step.Groups[j].Duration
		})
	}
	sort.Slice(slowSteps, func(i, j int) bool {
		return slowSteps[i].Duration > slowSteps[j].Duration
	})
	return job, slowSteps, nil
}

func getSlowSteps(steps []*github.TaskStep, threshold time.Duration) []*Step {
	slowSteps := make([]*Step, 0, len(steps))
	for _, s := range steps {
		step := &Step{
			Name:      s.GetName(),
			StartTime: *s.StartedAt.GetTime(),
			EndTime:   *s.CompletedAt.GetTime(),
		}
		step.Duration = step.EndTime.Sub(step.StartTime)
		if step.Duration < threshold {
			continue
		}
		slowSteps = append(slowSteps, step)
	}
	return slowSteps
}

func getSlowGroups(groups []*Group, threshold time.Duration) []*Group {
	slowGroups := make([]*Group, 0, len(groups))
	for _, group := range groups {
		if group.Duration < threshold {
			continue
		}
		slowGroups = append(slowGroups, group)
	}
	return slowGroups
}
