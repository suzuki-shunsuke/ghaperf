package view

import (
	"io"
	"sort"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
)

type Viewer struct {
	stdout io.Writer
}

func New(stdout io.Writer) *Viewer {
	return &Viewer{
		stdout: stdout,
	}
}

func jobDuration(job *github.WorkflowJob) time.Duration {
	return job.GetCompletedAt().Sub(job.GetStartedAt().Time)
}

type Step struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	duration  time.Duration
	Groups    []*parser.Group `json:"groups"`
}

func (s *Step) Duration() time.Duration {
	if s == nil {
		return 0
	}
	if s.EndTime.IsZero() {
		return 0
	}
	if s.duration != 0 {
		return s.duration
	}
	s.duration = s.EndTime.Sub(s.StartTime)
	return s.duration
}

func (s *Step) Contain(group *parser.Group, threshold time.Duration) {
	if group.Duration() < threshold {
		// The group is not slow
		// Go to the next group
		return
	}
	centerTime := s.StartTime.Add(s.Duration() / 2) //nolint:mnd
	if group.StartTime.Before(centerTime) && group.EndTime.After(centerTime) {
		// The group is contained in the step
		s.Groups = append(s.Groups, group)
		return
	}
}

func getSlowSteps(steps []*github.TaskStep, threshold time.Duration) []*Step {
	slowSteps := make([]*Step, 0, len(steps))
	for _, s := range steps {
		step := &Step{
			Name:      s.GetName(),
			StartTime: s.StartedAt.Time,
			EndTime:   s.CompletedAt.Time,
		}
		if step.Duration() < threshold {
			continue
		}
		slowSteps = append(slowSteps, step)
	}
	return slowSteps
}

func getSlowGroups(groups []*parser.Group, threshold time.Duration) []*parser.Group {
	slowGroups := make([]*parser.Group, 0, len(groups))
	for _, group := range groups {
		if group.Duration() < threshold {
			continue
		}
		slowGroups = append(slowGroups, group)
	}
	sort.Slice(slowGroups, func(i, j int) bool {
		return slowGroups[i].Duration() > slowGroups[j].Duration()
	})
	return slowGroups
}
