package view

import (
	"fmt"
	"sort"
	"time"

	"github.com/suzuki-shunsuke/ghaperf/pkg/parser"
)

func (v *Viewer) ShowGroups(groups []*parser.Group, threshold time.Duration) {
	slowGroups := getSlowGroups(groups, threshold)
	if len(slowGroups) == 0 {
		fmt.Fprintln(v.stdout, "No slow log group is found")
		return
	}
	sort.Slice(slowGroups, func(i, j int) bool {
		return slowGroups[i].Duration() > slowGroups[j].Duration()
	})

	fmt.Fprintln(v.stdout, "## Slow log groups")
	for i, group := range slowGroups {
		fmt.Fprintf(v.stdout, "%d. %s: %s\n", i+1, group.Duration(), group.Name)
	}
}
