package view

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v79/github"
	"github.com/suzuki-shunsuke/ghaperf/pkg/config"
)

type HeaderArg struct {
	Version                 string
	Repo                    string
	Now                     time.Time
	Threshold               time.Duration
	Count                   int
	WorkflowName            string
	ListWorkflowRunsOptions *github.ListWorkflowRunsOptions
	Config                  *config.Config
}

func (v *Viewer) ShowHeader(arg *HeaderArg) {
	fmt.Fprintln(v.stdout, generatedByText)
	var version string
	if arg.Version == "" || arg.Version == unknownVersion {
		version = unknownVersion
	} else {
		version = fmt.Sprintf(`<a href="https://github.com/suzuki-shunsuke/ghaperf/releases/tag/%s">%s</a>`, arg.Version, arg.Version)
	}
	fmt.Fprintln(v.stdout, "<table>")
	fmt.Fprintf(v.stdout, "<tr><td>ghaperf version</td><td>%s</td></tr>\n", version)
	fmt.Fprintf(v.stdout, "<tr><td>Created At</td><td>%s</td></tr>\n", arg.Now.Format(time.RFC3339))
	fmt.Fprintf(v.stdout, "<tr><td>Threshold</td><td>%s</td></tr>\n", arg.Threshold.Round(time.Second))
	fmt.Fprintf(v.stdout, `<tr><td>Repository</td><td><a href="https://github.com/%s">%s</a></td></tr>`+"\n", arg.Repo, arg.Repo)
	v.ShowConfigJobNames(arg)
	v.ShowConfigExcludedJobNames(arg)
	v.ShowConfigJobNameMappings(arg)
	if arg.Count > 0 {
		fmt.Fprintf(v.stdout, "<tr><td>The Number of Workflow Runs</td><td>%d</td></tr>\n", arg.Count)
	}
	if arg.WorkflowName != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Name</td><td>%s</td></tr>\n", arg.WorkflowName)
	}
	v.ShowListWorkflowRunsOptions(arg)
	fmt.Fprintf(v.stdout, "</table>\n\n")
}

func (v *Viewer) ShowConfigJobNames(arg *HeaderArg) {
	if arg.Config == nil || len(arg.Config.JobNames) == 0 {
		return
	}
	if len(arg.Config.JobNames) == 1 {
		fmt.Fprintf(v.stdout, "<tr><td>Job Names</td><td>%s</td></tr>\n", arg.Config.JobNames[0].String())
		return
	}
	names := make([]string, len(arg.Config.JobNames))
	for i, name := range arg.Config.JobNames {
		names[i] = fmt.Sprintf("<li>%s</li>", name.String())
	}
	fmt.Fprintf(v.stdout, "<tr><td>Job Names</td><td><ul>%s</ul></td></tr>\n", strings.Join(names, ""))
}

func (v *Viewer) ShowConfigExcludedJobNames(arg *HeaderArg) {
	if arg.Config == nil || len(arg.Config.ExcludedJobNames) == 0 {
		return
	}
	if len(arg.Config.ExcludedJobNames) == 1 {
		fmt.Fprintf(v.stdout, "<tr><td>Excluded Job Names</td><td>%s</td></tr>\n", arg.Config.ExcludedJobNames[0].String())
		return
	}
	names := make([]string, len(arg.Config.ExcludedJobNames))
	for i, name := range arg.Config.ExcludedJobNames {
		names[i] = fmt.Sprintf("<li>%s</li>", name.String())
	}
	fmt.Fprintf(v.stdout, "<tr><td>Excluded Job Names</td><td><ul>%s</ul></td></tr>\n", strings.Join(names, ""))
}

func (v *Viewer) ShowConfigJobNameMappings(arg *HeaderArg) {
	if arg.Config == nil || len(arg.Config.JobNameMappings) == 0 {
		return
	}
	if len(arg.Config.JobNameMappings) == 1 {
		for k, value := range arg.Config.JobNameMappings {
			fmt.Fprintf(v.stdout, "<tr><td>Job Name Mappings</td><td>%s => %s</td></tr>\n", k.String(), value)
		}
		return
	}
	names := make([]string, 0, len(arg.Config.JobNameMappings))
	for k, v := range arg.Config.JobNameMappings {
		names = append(names, fmt.Sprintf("<li>%s => %s</li>", k.String(), v))
	}
	fmt.Fprintf(v.stdout, "<tr><td>Job Name Mappings</td><td><ul>%s</ul></td></tr>\n", strings.Join(names, ""))
}

func (v *Viewer) ShowListWorkflowRunsOptions(arg *HeaderArg) {
	if arg.ListWorkflowRunsOptions == nil {
		return
	}
	if arg.ListWorkflowRunsOptions.Status != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Status</td><td>%s</td></tr>\n", arg.ListWorkflowRunsOptions.Status)
	}
	if arg.ListWorkflowRunsOptions.Actor != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Actor</td><td>%s</td></tr>\n", arg.ListWorkflowRunsOptions.Actor)
	}
	if arg.ListWorkflowRunsOptions.Branch != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Branch</td><td>%s</td></tr>\n", arg.ListWorkflowRunsOptions.Branch)
	}
	if arg.ListWorkflowRunsOptions.Event != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Event</td><td>%s</td></tr>\n", arg.ListWorkflowRunsOptions.Event)
	}
	if arg.ListWorkflowRunsOptions.Created != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Created</td><td>%s</td></tr>\n", arg.ListWorkflowRunsOptions.Created)
	}
}
