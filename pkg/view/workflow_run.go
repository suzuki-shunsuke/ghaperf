package view

import (
	"fmt"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
)

type JobWithSteps struct {
	Job       *github.WorkflowJob
	SlowSteps []*Step
	Duration  time.Duration
}

func (v *Viewer) ShowJobs(jobs []*github.WorkflowJob, threshold time.Duration) {
	arr := make([]*JobWithSteps, 0, len(jobs))
	for _, job := range jobs {
		d := jobDuration(job)
		if d < threshold {
			continue
		}
		slowSteps := getSlowSteps(job.Steps, threshold)
		arr = append(arr, &JobWithSteps{
			Job:       job,
			SlowSteps: slowSteps,
			Duration:  d,
		})
	}
	if len(arr) == 0 {
		fmt.Fprintln(v.stdout, "No slow jobs found")
		return
	}
	fmt.Fprintln(v.stdout, "## Slow jobs")
	for _, job := range arr {
		fmt.Fprintf(v.stdout, "Job Name: %s\n", job.Job.GetName())
		fmt.Fprintf(v.stdout, "Job ID: %d\n", job.Job.GetID())
		fmt.Fprintf(v.stdout, "Job Status: %s\n\n", job.Job.GetStatus())
		fmt.Fprintf(v.stdout, "Job Duration: %s\n\n", job.Duration)
		fmt.Fprintln(v.stdout, "### Slow steps")
		if len(job.SlowSteps) == 0 {
			fmt.Fprintf(v.stdout, "The job %s has no slow steps\n", job.Job.GetName())
			continue
		}
		for i, step := range job.SlowSteps {
			fmt.Fprintf(v.stdout, "%d. %s: %s\n", i+1, step.Duration(), step.Name)
			for j, group := range step.Groups {
				fmt.Fprintf(v.stdout, "   %d. %s: %s\n", j+1, group.Duration().Round(time.Second), group.Name)
			}
		}
	}
}
