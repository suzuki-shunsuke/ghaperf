# ghaperf

ghaperf is a CLI to analyze the performance of GitHub Actions using GitHub API and raw job logs

```sh
ghaperf --repo szksh-lab-2/test-github-action --job-id 53656655343 --threshold 1s
```

```markdown
Job Name: test
Job ID: 53656655343
Job Status: completed
Job Duration: 6s

## Slow steps
1. 2s: install aqua
   1. 1s: Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then
2. 2s: Run sleep 2
   1. 2s: Run sleep 2
3. 1s: Set up job
```

## :warning: This Project Is Still Work In Progress

This project is still work in progress.
Probably CLI doesn't work yet, and the document may be wrong.

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
- GITHUB_RUN_ID
- GITHUB_REPOSITORY
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
