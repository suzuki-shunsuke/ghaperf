package view

import (
	"fmt"
	"sort"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
)

type JobWithSteps struct {
	Job       *collector.Job
	SlowSteps []*Step
	Duration  time.Duration
}

func (v *Viewer) ShowJobs(run *github.WorkflowRun, jobs []*collector.Job, threshold time.Duration) {
	arr := make([]*JobWithSteps, 0, len(jobs))
	for _, job := range jobs {
		d := jobDuration(job.Job)
		if d < threshold {
			continue
		}
		arr = append(arr, &JobWithSteps{
			Job:      job,
			Duration: d,
		})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Duration > arr[j].Duration
	})
	fmt.Fprintf(v.stdout, "Workflow Run Name: %s\n", run.GetName())
	fmt.Fprintf(v.stdout, "Workflow Run ID: %d\n", run.GetID())
	fmt.Fprintf(v.stdout, "Workflow Run Status: %s\n", run.GetStatus())
	fmt.Fprintf(v.stdout, "Workflow Run Conclusion: %s\n", run.GetConclusion())
	fmt.Fprintf(v.stdout, "Workflow Run URL: %s\n", run.GetHTMLURL())
	for _, job := range arr {
		v.ShowJob(job.Job, threshold)
	}
}
