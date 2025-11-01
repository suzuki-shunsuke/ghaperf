package view

import (
	"fmt"
	"sort"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
)

type JobWithSteps struct {
	Job       *collector.Job
	SlowSteps []*Step
	Duration  time.Duration
}

func (v *Viewer) ShowJobs(run *collector.WorkflowRun, threshold time.Duration) {
	arr := make([]*JobWithSteps, 0, len(run.Jobs))
	for _, job := range run.Jobs {
		if job.Job.GetStatus() != "completed" {
			continue
		}
		d := job.Duration()
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
	fmt.Fprintln(v.stdout, "<table>")
	fmt.Fprintf(v.stdout, `<tr><td>Workflow Run Name</td><td><a href="%s">%s</a></td></tr>`+"\n", run.Run.GetHTMLURL(), run.Run.GetName())
	fmt.Fprintf(v.stdout, "<tr><td>Workflow Run ID</td><td>%d</td></tr>\n", run.Run.GetID())
	fmt.Fprintf(v.stdout, "<tr><td>Workflow Run Status</td><td>%s</td></tr>\n", run.Run.GetStatus())
	fmt.Fprintf(v.stdout, "<tr><td>Workflow Run Conclusion</td><td>%s</td></tr>\n", run.Run.GetConclusion())
	fmt.Fprintf(v.stdout, "</table>\n\n")
	if run.LogHasGone {
		v.ShowLogHasGone()
	}
	for _, job := range arr {
		v.ShowJob(job.Job, threshold)
	}
}
