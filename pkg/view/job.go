package view

import (
	"fmt"
	"sort"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
)

func (v *Viewer) ShowJob(j *collector.Job, threshold time.Duration) {
	job := j.Job
	slowSteps := getSlowSteps(job.Steps, threshold)
	sort.Slice(slowSteps, func(i, j int) bool {
		return slowSteps[i].Duration() > slowSteps[j].Duration()
	})

	allStepsDuration := time.Duration(0)
	for _, step := range job.Steps {
		allStepsDuration += step.GetCompletedAt().Sub(step.GetStartedAt().Time)
	}

	slowGroups := getSlowGroups(j.Groups, threshold)

	for _, step := range slowSteps {
		for _, group := range slowGroups {
			step.Contain(group, threshold)
		}
		sort.Slice(step.Groups, func(i, j int) bool {
			return step.Groups[i].Duration() > step.Groups[j].Duration()
		})
	}

	firstStepStartedAt := job.Steps[0].GetStartedAt().Time
	lastStepCompletedAt := job.Steps[len(job.Steps)-1].GetCompletedAt().Time

	fmt.Fprintf(v.stdout, "## Job: %s\n", job.GetName())

	fmt.Fprintln(v.stdout, "<table>")
	fmt.Fprintf(v.stdout, `<tr><td>Job ID</td><td><a href="%s">%d</a></td></tr>`+"\n", job.GetHTMLURL(), job.GetID())
	fmt.Fprintf(v.stdout, "<tr><td>Job Status</td><td>%s</td></tr>\n", job.GetStatus())
	fmt.Fprintf(v.stdout, "<tr><td>Job Conclusion</td><td>%s</td></tr>\n", job.GetConclusion())
	fmt.Fprintf(v.stdout, "<tr><td>Job Duration</td><td>%s</td></tr>\n", j.Duration())
	fmt.Fprintf(v.stdout, "<tr><td>All Steps Duration</td><td>%s</td></tr>\n", allStepsDuration.Round(time.Second))
	fmt.Fprintf(v.stdout, "<tr><td>Setup Job Duration</td><td>%s</td></tr>\n", firstStepStartedAt.Sub(job.StartedAt.Time).Round(time.Second))
	fmt.Fprintf(v.stdout, "<tr><td>Cleanup Job Duration</td><td>%s</td></tr>\n", job.GetCompletedAt().Sub(lastStepCompletedAt).Round(time.Second))
	fmt.Fprintf(v.stdout, "<tr><td>Steps Overhead</td><td>%s</td></tr>\n", (lastStepCompletedAt.Sub(firstStepStartedAt) - allStepsDuration).Round(time.Second))
	fmt.Fprintf(v.stdout, "</table>\n\n")

	if len(slowSteps) == 0 {
		fmt.Fprintf(v.stdout, "The job %s has no slow steps\n", job.GetName())
		return
	}

	fmt.Fprintln(v.stdout, "### Slow steps")
	for i, step := range slowSteps {
		fmt.Fprintf(v.stdout, "%d. %s: %s\n", i+1, step.Duration(), step.Name)
		if len(step.Groups) == 1 {
			continue
		}
		for j, group := range step.Groups {
			fmt.Fprintf(v.stdout, "   %d. %s: %s\n", j+1, group.Duration().Round(time.Second), group.Name)
		}
	}
}
