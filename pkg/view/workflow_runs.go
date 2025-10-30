package view

import (
	"fmt"
	"maps"
	"slices"
	"sort"
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

type JobMetric struct {
	Name   string
	Metric *Metric
	Steps  map[string]*StepMetric
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
		for _, job := range run.Jobs {
			if job.Job.GetStatus() != "completed" {
				continue
			}
			// TODO normalize job name for matrix jobs
			jobName := job.Job.GetName()
			jm, ok := jobMetrics[jobName]
			if !ok {
				jm = &JobMetric{
					Name:   jobName,
					Metric: &Metric{},
					Steps:  map[string]*StepMetric{},
				}
				jobMetrics[jobName] = jm
			}
			jm.Metric.Add(job.Duration())

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
		fmt.Fprintf(v.stdout, "Total Job Duration: %s\n", jm.Metric.Sum.Round(time.Second))
		fmt.Fprintf(v.stdout, "The number of Job Executions: %d\n", jm.Metric.Count)
		fmt.Fprintf(v.stdout, "Average Job Duration: %s\n", jm.Metric.Avg.Round(time.Second))
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
			fmt.Fprintf(v.stdout, "%d. %s: total:%s, count:%d, avg:%s\n", i+1, sm.Name, sm.Metric.Sum, sm.Metric.Count, sm.Metric.Avg.Round(time.Second))
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
				fmt.Fprintf(v.stdout, "    %d. %s: total:%s, count:%d, avg:%s\n", j+1, gm.Name, gm.Metric.Sum.Round(time.Second), gm.Metric.Count, gm.Metric.Avg.Round(time.Second))
			}
		}
	}
}
