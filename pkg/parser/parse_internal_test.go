package parser

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_parseLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		txt  string
		line *Line
	}{
		{
			name: "invalid",
			txt:  "invalid",
			line: &Line{
				Continue: true,
				Content:  "invalid",
			},
		},
		{
			name: "invalid timestamp",
			txt:  "hello world",
			line: &Line{
				Continue: true,
				Content:  "hello world",
			},
		},
		{
			name: "##[group]",
			txt:  "2025-10-25T13:48:59.4421674Z ##[group]Runner Image Provisioner",
			line: &Line{
				Content:   "Runner Image Provisioner",
				Start:     true,
				Timestamp: time.Date(2025, 10, 25, 13, 48, 59, 4421674, time.UTC),
			},
		},
		{
			name: "job name",
			txt:  "2025-10-29T13:56:22.7273757Z Complete job name: test / test / test (windows-latest, arm64)",
			line: &Line{
				Content:   "Complete job name: test / test / test (windows-latest, arm64)",
				JobName:   "test / test / test (windows-latest, arm64)",
				Timestamp: time.Date(2025, 10, 29, 13, 56, 22, 7273757, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			line := parseLine(tt.txt)
			if diff := cmp.Diff(tt.line, line, cmp.AllowUnexported(Line{})); diff != "" {
				t.Errorf("line parseLogLine() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
