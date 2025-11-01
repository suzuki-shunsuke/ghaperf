package view

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/collector"
	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
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

func (v *Viewer) ShowRuns(runs []*collector.WorkflowRun, threshold time.Duration) {
	jobMetrics := map[string]*JobMetric{}
	for _, run := range runs {
		setMetricsByRun(jobMetrics, run)
	}
	jobArr := slices.Collect(maps.Values(jobMetrics))
	// extract only slow jobs
	slowJobs := getSlowJobs(jobArr, threshold)
	if len(slowJobs) == 0 {
		fmt.Fprintln(v.stdout, "There is no slow job")
		return
	}
	for _, jm := range slowJobs {
		v.showJobMetric(jm, threshold)
	}
}

func (v *Viewer) showJobMetric(jm *JobMetric, threshold time.Duration) {
	if jm.Metric.Avg < threshold {
		return
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
		return
	}
	v.showSlowStepMetrics(slowSteps, threshold)
}

func (v *Viewer) showSlowStepMetrics(slowSteps []*StepMetric, threshold time.Duration) {
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

func getSlowJobs(jobs []*JobMetric, threshold time.Duration) []*JobMetric {
	arr := make([]*JobMetric, 0, len(jobs))
	for _, jm := range jobs {
		if jm.Metric.Avg < threshold {
			continue
		}
		arr = append(arr, jm)
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Metric.Sum > arr[j].Metric.Sum
	})
	return arr
}

func setMetricsByRun(jobMetrics map[string]*JobMetric, run *collector.WorkflowRun) {
	slowestJobs := map[string]*collector.Job{}
	for _, job := range run.Jobs {
		setMetricsByJob(jobMetrics, slowestJobs, job)
	}
	for normalizedJobName, job := range slowestJobs {
		jobMetrics[normalizedJobName].Metric.Add(job.Duration())
	}
}

func setMetricsByJob(jobMetrics map[string]*JobMetric, slowestJobs map[string]*collector.Job, job *collector.Job) {
	if job.Job.GetStatus() != "completed" {
		return
	}
	// Get the slowest job for each normalized job name
	if slowestJobs[job.NormalizedName].Duration() < job.Duration() {
		slowestJobs[job.NormalizedName] = job
	}
	jm := initJobMetric(jobMetrics, job)
	updateSlowestJobs(jm, job)

	for _, s := range job.Job.Steps {
		setJobMetric(jm, job, s)
	}
}

func setJobMetric(jm *JobMetric, job *collector.Job, s *github.TaskStep) {
	sm := initStepMetric(jm, s)
	step := &Step{
		Name:      s.GetName(),
		StartTime: s.StartedAt.Time,
		EndTime:   s.CompletedAt.Time,
	}
	sm.Metric.Add(step.Duration())
	// Extract groups belonging to the step
	for _, group := range job.Groups {
		step.Contain(group)
	}
	// Add step groups to the step metric
	for _, group := range step.Groups {
		m := initGroupMetric(sm, group)
		m.Add(group.Duration())
	}
}

func initGroupMetric(sm *StepMetric, group *parser.Group) *Metric {
	m, ok := sm.Groups[group.Name]
	if ok {
		return m
	}
	m = &Metric{}
	sm.Groups[group.Name] = m
	return m
}

func initStepMetric(jm *JobMetric, step *github.TaskStep) *StepMetric {
	sm, ok := jm.Steps[step.GetName()]
	if ok {
		return sm
	}
	sm = &StepMetric{
		Name:   step.GetName(),
		Metric: &Metric{},
		Groups: map[string]*Metric{},
	}
	jm.Steps[step.GetName()] = sm
	return sm
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

func initJobMetric(jobMetrics map[string]*JobMetric, job *collector.Job) *JobMetric {
	jm, ok := jobMetrics[job.NormalizedName]
	if !ok {
		// Initialize JobMetric
		jm = &JobMetric{
			Name:        job.NormalizedName,
			Metric:      &Metric{},
			Steps:       map[string]*StepMetric{},
			SlowestJobs: make([]*collector.Job, 0, countSlowest),
		}
		jobMetrics[job.NormalizedName] = jm
	}
	return jm
}
