# ghaperf

ghaperf is a CLI to analyze the performance of GitHub Actions using GitHub API and raw job logs.
It's useful to detect the bottlenecks inside composite actions.
It outputs a report about the performance with the markdown format.

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

[There are some awesome tools](#similar-works), but they can't retrieve step-level data inside composite actions because [the Workflow Jobs APIs](https://docs.github.com/en/rest/actions/workflow-jobs) don’t include those data.
On the other hand, ghaperf can detect bottlenecks within composite actions.
It retrieves job logs via the API, parses them, and extracts data from all log groups — including steps within composite actions.

Note that the specification of log format isn't published, so ghaperf may fail to parse logs due to unexpected specification, and it may get broken suddenly due to changes of the log specification.

## Install

```sh
go install github.com/suzuki-shunsuke/ghaperf/cmd/ghaperf@latest
```

## Getting Started

[A GitHub Access token is required to avoid API rate limit or to access private repositories.](#github-access-token)

```sh
export GITHUB_TOKEN=xxx
```

1. Run against a log file ([example](https://github.com/suzuki-shunsuke/ghaperf/blob/main/testdata/log.txt)):

```sh
git clone https://github.com/suzuki-shunsuke/ghaperf
cd ghaperf
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

ghaperf outputs the report.
ghaperf reports steps and log groups which take longer than threshold.

e.g.

```markdown
### Slow steps
1. 2s: install aqua
   1. 1s: Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then
2. 2s: Run sleep 2
   1. 2s: Run sleep 2
3. 1s: Set up job
```

The default threshold is `30s`, but you can change this by `--threshold` option and the environment variable `GHAPERF_THRESHOLD`.
It's parsed by [time.ParseDuration](https://pkg.go.dev/time#ParseDuration).

## Environment Variables

- GHAPERF_LOG_LEVEL: `debug|info|warn|error`. The default is `info`
- GHAPERF_GITHUB_TOKEN
- GHAPERF_GHTKN
- GHAPERF_THRESHOLD: The threshold of steps and log groups' duration
- GITHUB_TOKEN

### GitHub Access Token

A GitHub access token is required to get workflow runs and jobs and their logs.
Private repositories require the `Actions: Read` permission.

```sh
export GITHUB_TOKEN=xxx
# or
export GHAPERF_GITHUB_TOKEN=xxx
```

Or if you use [ghtkn](https://github.com/suzuki-shunsuke/ghtkn), you can enable the integration.

```sh
export GHAPERF_GHTKN=true
```

## Cache

ghaperf caches raw data of completed workflow runs and jobs in the cache directory `${XDG_CACHE_HOME:-${HOME}/.cache}/ghaperf/`.

## :warning: Note

ghaperf gets job logs by [GitHub API](https://docs.github.com/en/rest/actions/workflow-jobs#download-job-logs-for-a-workflow-run), but if jobs aren't completed or completed just now, the API would fail.

## Similar Works

- https://github.com/Kesin11/actions-timeline
- https://github.com/Kesin11/github_actions_otel_trace
- https://github.com/Kesin11/CIAnalyzer
- https://github.com/inception-health/otel-export-trace-action
- https://github.com/runforesight/workflow-telemetry-action
- https://github.com/paper2/github-actions-opentelemetry
