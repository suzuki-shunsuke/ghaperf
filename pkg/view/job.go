package view

import (
	"fmt"
	"sort"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
)

func (v *Viewer) ShowJob(groups []*parser.Group, threshold time.Duration, job *github.WorkflowJob) {
	slowSteps := getSlowSteps(job.Steps, threshold)
	if len(slowSteps) == 0 {
		fmt.Fprintln(v.stdout, "No slow step is found")
		return
	}
	sort.Slice(slowSteps, func(i, j int) bool {
		return slowSteps[i].Duration() > slowSteps[j].Duration()
	})

	slowGroups := getSlowGroups(groups, threshold)
	sort.Slice(slowGroups, func(i, j int) bool {
		return slowGroups[i].Duration() > slowGroups[j].Duration()
	})

	fmt.Fprintln(v.stdout, "## Slow log groups")
	for i, group := range slowGroups {
		fmt.Fprintf(v.stdout, "%d. %s: %s\n", i+1, group.Duration(), group.Name)
	}

	fmt.Fprintln(v.stdout, "## Slow steps")
	fmt.Fprintf(v.stdout, "Job Name: %s\n", job.GetName())
	fmt.Fprintf(v.stdout, "Job ID: %d\n", job.GetID())
	fmt.Fprintf(v.stdout, "Job Status: %s\n\n", job.GetStatus())
	fmt.Fprintf(v.stdout, "Job Duration: %s\n\n", jobDuration(job))
	if len(slowSteps) == 0 {
		fmt.Fprintf(v.stdout, "The job %s has no slow steps\n", job.GetName())
		return
	}
	for i, step := range slowSteps {
		fmt.Fprintf(v.stdout, "%d. %s: %s\n", i+1, step.Duration(), step.Name)
		for j, group := range step.Groups {
			fmt.Fprintf(v.stdout, "   %d %s: %s\n", j+1, group.Duration(), group.Name)
		}
	}
}
