package view

import (
	"fmt"
	"time"

	"github.com/google/go-github/v76/github"
)

type HeaderArg struct {
	Version                 string
	Now                     time.Time
	Threshold               time.Duration
	Count                   int
	WorkflowName            string
	ListWorkflowRunsOptions *github.ListWorkflowRunsOptions
}

func (v *Viewer) ShowHeader(arg *HeaderArg) { //nolint:cyclop
	fmt.Fprintln(v.stdout, generatedByText)
	var version string
	if arg.Version == "" || arg.Version == unknownVersion {
		version = unknownVersion
	} else {
		version = fmt.Sprintf(`<a href="https://github.com/suzuki-shunsuke/ghaperf/releases/tag/v%s">v%s</a>`, arg.Version, arg.Version)
	}
	fmt.Fprintln(v.stdout, "<table>")
	fmt.Fprintf(v.stdout, `<tr><td>ghaperf version</td><td>%s</td></tr>`+"\n", version)
	fmt.Fprintf(v.stdout, "<tr><td>Created At</td><td>%s</td></tr>\n", arg.Now.Format(time.RFC3339))
	fmt.Fprintf(v.stdout, "<tr><td>Threshold</td><td>%s</td></tr>\n", arg.Threshold.Round(time.Second))
	if arg.Count > 0 {
		fmt.Fprintf(v.stdout, "<tr><td>The Number of Workflow Runs</td><td>%d</td></tr>\n", arg.Count)
	}
	if arg.WorkflowName != "" {
		fmt.Fprintf(v.stdout, "<tr><td>Workflow Name</td><td>%s</td></tr>\n", arg.WorkflowName)
	}
	if arg.ListWorkflowRunsOptions != nil { //nolint:nestif
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
	fmt.Fprintf(v.stdout, "</table>\n\n")
}
