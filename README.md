# ghaperf

ghaperf is a CLI to analyze the performance of GitHub Actions using GitHub API and raw job logs.

```sh
ghaperf --repo szksh-lab-2/test-github-action --job-id 53656655343 --threshold 1s
```

```markdown
## Job: test
Job Name: test
Job ID: 53656655343
Job URL: https://github.com/szksh-lab-2/test-github-action/actions/runs/18804539465/job/53656655343
Job Status: completed
Job Conclusion: success
Job Duration: 6s
All Steps Duration: 5s
Setup Job Duration: 0s
Cleanup Job Duration: 1s
Steps Overhead: 0s

### Slow steps
1. 2s: install aqua
   1. 1s: Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then
2. 2s: Run sleep 2
   1. 2s: Run sleep 2
3. 1s: Set up job
```

## :warning: This Project Is Still Work In Progress

This project is still work in progress.
Probably CLI doesn't work yet, and the document may be wrong.

## Why?

Unlike [other tools](#similar-works), ghaperf can detect bottlenecks within composite actions.
Other tools use [the Workflow Jobs APIs](https://docs.github.com/en/rest/actions/workflow-jobs) to get workflow runs and jobs data, but these APIs don’t include step-level data inside composite actions.
As a result, even if you can identify a slow composite action, you can’t tell which specific steps in the action are causing the slowdown.
To address this limitation, ghaperf retrieves job logs via the API, parses them, and extracts data from all log groups — including steps within composite actions.

## Install

```
go install github.com/suzuki-shunsuke/ghaperf/cmd/ghaperf@latest
```

## GitHub Access Token

A GitHub access token is required to get workflow runs and jobs and their logs.
Private repositories require the `Actions: Read` permission.

## Environment Variables

- GHAPERF_LOG_LEVEL
- GHAPERF_GITHUB_TOKEN
- GHAPERF_GHTKN
- GHAPERF_THRESHOLD
- GITHUB_TOKEN

## Getting Started

1. Run against a log file:

```sh
ghaperf --log-file testdata/log.txt
```

2. Run against a job:

```sh
ghaperf --repo szksh-lab-2/test-github-action --job-id 53656655343 --threshold 1s
```

3. Run against a workflow run:

```sh
ghaperf --repo szksh-lab-2/test-github-action --run-id 18804539465 --threshold 1s
```

## :bulb: Cache

ghaperf caches raw data of completed workflow runs and jobs in the cache directory `${XDG_CACHE_HOME:-${HOME}/.cache}/ghaperf/`.

## Similar Works

- https://github.com/Kesin11/actions-timeline
- https://github.com/Kesin11/github_actions_otel_trace
- https://github.com/Kesin11/CIAnalyzer
- https://github.com/inception-health/otel-export-trace-action
- https://github.com/runforesight/workflow-telemetry-action
- https://github.com/paper2/github-actions-opentelemetry
