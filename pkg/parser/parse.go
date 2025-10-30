package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

/*
2025-10-25T13:48:59.4397988Z Current runner version: '2.329.0'
2025-10-25T13:48:59.4421674Z ##[group]Runner Image Provisioner
2025-10-25T13:48:59.4422447Z Hosted Compute Agent
2025-10-25T13:48:59.4422930Z Version: 20250912.392
2025-10-25T13:48:59.4423518Z Commit: d921fda672a98b64f4f82364647e2f10b2267d0b
2025-10-25T13:48:59.4424160Z Build Date: 2025-09-12T15:23:14Z
2025-10-25T13:48:59.4424674Z ##[endgroup]
2025-10-25T13:48:59.4425179Z ##[group]Operating System
2025-10-25T13:48:59.4425653Z Ubuntu
2025-10-25T13:48:59.4426139Z 24.04.3
2025-10-25T13:48:59.4426535Z LTS
2025-10-25T13:48:59.4426940Z ##[endgroup]
2025-10-25T13:48:59.4427390Z ##[group]Runner Image
2025-10-25T13:48:59.4427903Z Image: ubuntu-24.04
2025-10-25T13:48:59.4428352Z Version: 20250929.60.1
2025-10-25T13:48:59.4429743Z Included Software: https://github.com/actions/runner-images/blob/ubuntu24/20250929.60/images/ubuntu/Ubuntu2404-Readme.md
2025-10-25T13:48:59.4431149Z Image Release: https://github.com/actions/runner-images/releases/tag/ubuntu24%2F20250929.60
2025-10-25T13:48:59.4432021Z ##[endgroup]
2025-10-25T13:48:59.4432847Z ##[group]GITHUB_TOKEN Permissions
2025-10-25T13:48:59.4434949Z Metadata: read
2025-10-25T13:48:59.4435405Z ##[endgroup]
2025-10-25T13:48:59.4437344Z Secret source: Actions
2025-10-25T13:48:59.4437953Z Prepare workflow directory
2025-10-25T13:48:59.4747004Z Prepare all required actions
2025-10-25T13:48:59.4838926Z Complete job name: test
2025-10-25T13:48:59.5432811Z ##[group]Run echo "start"
2025-10-25T13:48:59.5433375Z [36;1mecho "start"[0m
2025-10-25T13:48:59.7049889Z shell: /usr/bin/bash -e {0}
2025-10-25T13:48:59.7051322Z ##[endgroup]
2025-10-25T13:48:59.7353615Z start
2025-10-25T13:48:59.7447602Z ##[group]Run sleep 2
2025-10-25T13:48:59.7448098Z [36;1msleep 2[0m
2025-10-25T13:48:59.7466996Z shell: /usr/bin/bash -e {0}
2025-10-25T13:48:59.7467478Z ##[endgroup]
2025-10-25T13:49:01.7606729Z ##[group]Run echo "end"
2025-10-25T13:49:01.7607934Z [36;1mecho "end"[0m
2025-10-25T13:49:01.7630244Z shell: /usr/bin/bash -e {0}
2025-10-25T13:49:01.7631575Z ##[endgroup]
2025-10-25T13:49:01.7683231Z end
2025-10-25T13:49:01.7764649Z Cleaning up orphan processes
*/

type Group struct {
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	duration  time.Duration
	Lines     []*Line `json:"lines"`
}

func (g *Group) Duration() time.Duration {
	if g == nil {
		return 0
	}
	if g.duration != 0 {
		return g.duration
	}
	if g.EndTime.IsZero() {
		return 0
	}
	g.duration = g.EndTime.Sub(g.StartTime)
	return g.duration
}

type Line struct {
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Start     bool
	Continue  bool
	JobName   string
}

type Log struct {
	JobName  string
	Groups   []*Group
	duration time.Duration
}

func (l *Log) Duration() time.Duration {
	if l == nil {
		return 0
	}
	if l.duration != 0 {
		return l.duration
	}
	if len(l.Groups) == 0 {
		return 0
	}
	l.duration = l.Groups[len(l.Groups)-1].EndTime.Sub(l.Groups[0].StartTime)
	return l.duration
}

func Parse(data io.Reader) (*Log, error) {
	scanner := bufio.NewScanner(data)
	log := &Log{}
	group := &Group{}

	for scanner.Scan() {
		txt := scanner.Text()
		line := parseLine(txt)
		if log.JobName == "" && line.JobName != "" {
			log.JobName = line.JobName
		}
		if line.Continue {
			group.Lines[len(group.Lines)-1].Content += "\n" + line.Content
			continue
		}
		group.Lines = append(group.Lines, line)
		if line.Start {
			// End the previous group
			group.EndTime = line.Timestamp
			log.Groups = append(log.Groups, group)
			group = &Group{}
		}
	}
	if group.EndTime.IsZero() {
		group.EndTime = group.Lines[len(group.Lines)-1].Timestamp
	}
	log.Groups = append(log.Groups, group)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan a log file: %w", err)
	}

	return log, nil
}

var ansiEscapeSequence = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

func parseLine(txt string) *Line {
	txt = ansiEscapeSequence.ReplaceAllString(
		strings.TrimPrefix(txt, "\ufeff"), "") // Remove BOM and ANSI escape sequences
	d, l, ok := strings.Cut(txt, " ")
	if !ok {
		// The log doesn't start with timestamp.
		// This is a continuation from the previous log.
		return &Line{
			Continue: true,
			Content:  txt,
		}
	}
	// 2025-10-25T13:48:59.4421674Z ##[group]Runner Image Provisioner
	t, err := time.Parse("2006-01-02T15:04:05.9999999Z", d)
	if err != nil {
		// The log doesn't start with timestamp.
		// This is a continuation from the previous log.
		return &Line{
			Continue: true,
			Content:  txt,
		}
	}

	switch {
	case strings.HasPrefix(l, "##[group]"):
		return &Line{
			Start:     true,
			Content:   strings.TrimPrefix(l, "##[group]"),
			Timestamp: t,
		}
	case strings.HasPrefix(l, "Complete job name: "):
		// 2025-10-29T13:56:22.7273757Z Complete job name: test / test / test (windows-latest, arm64)
		return &Line{
			Content:   l,
			JobName:   strings.TrimPrefix(l, "Complete job name: "),
			Timestamp: t,
		}
	default:
		return &Line{
			Content:   l,
			Timestamp: t,
		}
	}
}
