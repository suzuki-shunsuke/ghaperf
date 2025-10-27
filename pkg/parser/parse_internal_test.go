//nolint:funlen
package parser

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_parseLogLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		line        string
		group       *Group
		newGroup    *Group
		passedGroup *Group
		checkErr    func(*testing.T, error)
	}{
		{
			name: "invalid",
			line: "invalid",
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if !errors.Is(err, errInvalidLogLineFormat) {
					t.Fatalf("expected errInvalidLogLineFormat but got: %v", err)
				}
			},
		},
		{
			name: "invalid timestamp",
			line: "hello world",
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				var e *time.ParseError
				if !errors.As(err, &e) {
					t.Fatalf("expected time.ParseError but got: %v", err)
				}
			},
		},
		{
			name: "##[group]",
			line: "2025-10-25T13:48:59.4421674Z ##[group]Runner Image Provisioner",
			newGroup: &Group{
				Name:      "Runner Image Provisioner",
				StartTime: time.Date(2025, 10, 25, 13, 48, 59, 442167400, time.UTC),
			},
		},
		{
			name: "##[group]",
			line: "2025-10-25T13:48:59.4425179Z ##[group]Operating System",
			group: &Group{
				Name:      "Runner Image Provisioner",
				StartTime: time.Date(2025, 10, 25, 13, 48, 59, 442167400, time.UTC),
			},
			passedGroup: &Group{
				Name:      "Runner Image Provisioner",
				StartTime: time.Date(2025, 10, 25, 13, 48, 59, 442167400, time.UTC),
				EndTime:   time.Date(2025, 10, 25, 13, 48, 59, 442517900, time.UTC),
			},
			newGroup: &Group{
				Name:      "Operating System",
				StartTime: time.Date(2025, 10, 25, 13, 48, 59, 442517900, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			group, err := parseLogLine(tt.line, tt.group)
			if err != nil {
				if tt.checkErr == nil {
					t.Fatalf("parseLogLine() unexpected error: %v", err)
				}
				tt.checkErr(t, err)
			}
			if diff := cmp.Diff(tt.newGroup, group, cmp.AllowUnexported(Group{})); diff != "" {
				t.Errorf("newGroup parseLogLine() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.passedGroup, tt.group, cmp.AllowUnexported(Group{})); diff != "" {
				t.Errorf("passedGroup parseLogLine() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
