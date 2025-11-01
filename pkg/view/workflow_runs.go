package view

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
)

type Metric struct {
	Sum   time.Duration
	Count int
	Avg   time.Duration
}

func (m *Metric) Add(d time.Duration) {
	m.Sum += d
	m.Count++
	m.Avg = m.Sum / time.Duration(m.Count)
}

const countSlowest = 3

type JobMetric struct {
	Name        string
	Metric      *Metric
	Steps       map[string]*StepMetric
	SlowestJobs []*collector.Job
}

type StepMetric struct {
	Name   string
	Metric *Metric
	Groups map[string]*Metric
}

type GroupMetric struct {
	Name   string
	Metric *Metric
}

func (v *Viewer) ShowRuns(runs []*collector.WorkflowRun, threshold time.Duration) { //nolint:gocognit,cyclop,funlen
	jobMetrics := map[string]*JobMetric{}
	for _, run := range runs {
		normalizedJobs := map[string]*collector.Job{}
		for _, job := range run.Jobs {
			if job.Job.GetStatus() != "completed" {
				continue
			}
			if normalizedJobs[job.NormalizedName].Duration() < job.Duration() {
				normalizedJobs[job.NormalizedName] = job
			}
			jm, ok := jobMetrics[job.NormalizedName]
			if !ok {
				jm = &JobMetric{
					Name:        job.NormalizedName,
					Metric:      &Metric{},
					Steps:       map[string]*StepMetric{},
					SlowestJobs: make([]*collector.Job, 0, countSlowest),
				}
				jobMetrics[job.NormalizedName] = jm
			}
			updateSlowestJobs(jm, job)

			for _, s := range job.Job.Steps {
				sm, ok := jm.Steps[s.GetName()]
				if !ok {
					sm = &StepMetric{
						Name:   s.GetName(),
						Metric: &Metric{},
						Groups: map[string]*Metric{},
					}
					jm.Steps[s.GetName()] = sm
				}
				step := &Step{
					Name:      s.GetName(),
					StartTime: s.StartedAt.Time,
					EndTime:   s.CompletedAt.Time,
				}
				sm.Metric.Add(step.Duration())
				for _, group := range job.Groups {
					step.Contain(group, threshold)
				}
				for _, group := range step.Groups {
					m, ok := sm.Groups[group.Name]
					if !ok {
						m = &Metric{}
						sm.Groups[group.Name] = m
					}
					m.Add(group.Duration())
				}
			}
		}
		for normalizedJobName, job := range normalizedJobs {
			jobMetrics[normalizedJobName].Metric.Add(job.Duration())
		}
	}
	jobArr := slices.Collect(maps.Values(jobMetrics))
	arr := make([]*JobMetric, 0, len(jobArr))
	for _, jm := range jobArr {
		if jm.Metric.Avg < threshold {
			continue
		}
		arr = append(arr, jm)
	}
	if len(arr) == 0 {
		fmt.Fprintln(v.stdout, "There is no slow job")
		return
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Metric.Sum > arr[j].Metric.Sum
	})
	for _, jm := range arr {
		if jm.Metric.Avg < threshold {
			continue
		}
		stepArr := slices.Collect(maps.Values(jm.Steps))
		sort.Slice(stepArr, func(i, j int) bool {
			return stepArr[i].Metric.Sum > stepArr[j].Metric.Sum
		})
		fmt.Fprintf(v.stdout, "## Job: %s\n", jm.Name)
		slowestJobStrs := make([]string, len(jm.SlowestJobs))
		for i, job := range jm.SlowestJobs {
			slowestJobStrs[i] = fmt.Sprintf(`<a href="%s">%s</a>`, job.Job.GetHTMLURL(), job.Duration())
		}
		fmt.Fprintln(v.stdout, "<table>")
		fmt.Fprintf(v.stdout, "<tr><td>Average Job Duration</td><td>%s (%s/%d)</td></tr>\n", jm.Metric.Avg.Round(time.Second), jm.Metric.Sum.Round(time.Second), jm.Metric.Count)
		fmt.Fprintf(v.stdout, "<tr><td>Slowest Jobs</td><td>%s</td></tr>\n", strings.Join(slowestJobStrs, ", "))
		fmt.Fprintf(v.stdout, "</table>\n\n")
		slowSteps := make([]*StepMetric, 0, len(stepArr))
		for _, sm := range stepArr {
			if sm.Metric.Avg < threshold {
				continue
			}
			slowSteps = append(slowSteps, sm)
		}
		if len(slowSteps) == 0 {
			fmt.Fprintln(v.stdout, "The job has no slow steps")
			continue
		}
		fmt.Fprintln(v.stdout, "### Slow steps")
		for i, sm := range slowSteps {
			fmt.Fprintf(v.stdout, "%d. %s (%s/%d): %s\n", i+1, sm.Metric.Avg.Round(time.Second), sm.Metric.Sum, sm.Metric.Count, sm.Name)
			if len(sm.Groups) <= 1 {
				continue
			}
			groupArr := make([]*GroupMetric, 0, len(sm.Groups))
			for groupName, m := range sm.Groups {
				groupArr = append(groupArr, &GroupMetric{
					Name:   groupName,
					Metric: m,
				})
			}
			sort.Slice(groupArr, func(i, j int) bool {
				return groupArr[i].Metric.Sum > groupArr[j].Metric.Sum
			})
			for j, gm := range groupArr {
				if gm.Metric.Avg < threshold {
					continue
				}
				fmt.Fprintf(v.stdout, "    %d. %s (%s/%d): %s\n", j+1, gm.Metric.Avg.Round(time.Second), gm.Metric.Sum.Round(time.Second), gm.Metric.Count, gm.Name)
			}
		}
	}
}

func updateSlowestJobs(jm *JobMetric, job *collector.Job) {
	if len(jm.SlowestJobs) < countSlowest {
		jm.SlowestJobs = append(jm.SlowestJobs, job)
		if len(jm.SlowestJobs) != countSlowest {
			return
		}
		sort.Slice(jm.SlowestJobs, func(i, j int) bool {
			return jm.SlowestJobs[i].Duration() > jm.SlowestJobs[j].Duration()
		})
		return
	}
	if jm.SlowestJobs[countSlowest-1].Duration() >= job.Duration() {
		return
	}
	jm.SlowestJobs[countSlowest-1] = job
	sort.Slice(jm.SlowestJobs, func(i, j int) bool {
		return jm.SlowestJobs[i].Duration() > jm.SlowestJobs[j].Duration()
	})
}
