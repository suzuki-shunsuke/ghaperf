# ghaperf

ghaperf is a CLI to analyze the performance of GitHub Actions using GitHub API and raw job logs.
It's useful to detect the bottlenecks inside composite actions.
It outputs a report about the performance with the markdown format.

1. Analyze a workflow's performance:

```sh
ghaperf --repo aquaproj/aqua-registry --workflow test.yaml --count 10 --threshold 2s
```

<details>
<summary>Generated Report</summary>

## Job: test / test / test (windows-latest)
- Total Job Duration: 6m3s
- The number of Job Executions: 10
- Average Job Duration: 36s
- Slowest jobs: [45s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227585996), [44s](https://github.com/aquaproj/aqua-registry/actions/runs/18984012981/job/54223408471), [38s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223445810)
### Slow steps
1. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 14s (2m21s/10)
    1. Run $(if($env:AQUA_ROOT_DIR) {echo $env:AQUA_ROOT_DIR} else {echo "$HOME/AppData/Local/aquaproj-aqua/bin"}) | Out-File -FilePath $env:GITHUB_PATH -Encoding utf8 -Append: 9s (1m28s/10)
    2. Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then: 5s (46s/10)
2. Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8: 8s (1m15s/10)
    1. Setting up auth: 2s (16s/7)
    2. Fetching the repository: 2s (5s/2)
3. Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8: 6s (56s/10)
    1. Run aqua i --test: 6s (34s/6)
    2. Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV": 2s (17s/7)
4. Set up job: 4s (37s/10)
## Job: test / test / test (windows-latest, arm64)
- Total Job Duration: 5m59s
- The number of Job Executions: 10
- Average Job Duration: 36s
- Slowest jobs: [51s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223445799), [40s](https://github.com/aquaproj/aqua-registry/actions/runs/18983245993/job/54220759743), [38s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227586013)
### Slow steps
1. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 14s (2m19s/10)
    1. Run $(if($env:AQUA_ROOT_DIR) {echo $env:AQUA_ROOT_DIR} else {echo "$HOME/AppData/Local/aquaproj-aqua/bin"}) | Out-File -FilePath $env:GITHUB_PATH -Encoding utf8 -Append: 9s (1m26s/10)
    2. Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then: 5s (47s/10)
2. Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8: 7s (1m14s/10)
    1. Setting up auth: 3s (16s/6)
    2. Fetching the repository: 3s (3s/1)
3. Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8: 6s (57s/10)
    1. Run aqua i --test: 7s (29s/4)
    2. Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV": 2s (14s/6)
4. Set up job: 4s (41s/10)
## Job: test / test / test (macos-13)
- Total Job Duration: 4m59s
- The number of Job Executions: 10
- Average Job Duration: 30s
- Slowest jobs: [39s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223445786), [36s](https://github.com/aquaproj/aqua-registry/actions/runs/18982948414/job/54219767595), [33s](https://github.com/aquaproj/aqua-registry/actions/runs/18985653020/job/54228560759)
### Slow steps
1. Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8: 7s (1m13s/10)
    1. Run aqua i --test: 7s (49s/7)
    2. Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV": 2s (18s/8)
2. Set up job: 7s (1m8s/10)
3. Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8: 5s (53s/10)
    1. Checking out the ref: 3s (20s/7)
    2. Fetching the repository: 3s (3s/1)
4. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 5s (45s/10)
## Job: test / test / test (macos-14)
- Total Job Duration: 4m45s
- The number of Job Executions: 10
- Average Job Duration: 29s
- Slowest jobs: [1m45s](https://github.com/aquaproj/aqua-registry/actions/runs/18985291460/job/54227488846), [23s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223445788), [22s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227585994)
### Slow steps
1. Set up job: 6s (1m0s/10)
2. Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8: 5s (45s/10)
    1. Run aqua i --test: 4s (25s/7)
    2. Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV": 2s (6s/3)
3. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 3s (30s/10)
4. Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8: 3s (26s/10)
## Job: test / lintnet / lintnet
- Total Job Duration: 2m39s
- The number of Job Executions: 10
- Average Job Duration: 16s
- Slowest jobs: [18s](https://github.com/aquaproj/aqua-registry/actions/runs/18985291460/job/54227470425), [18s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223426855), [18s](https://github.com/aquaproj/aqua-registry/actions/runs/18982975753/job/54219833548)
### Slow steps
1. Run lintnet lint: 6s (1m2s/10)
2. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 3s (29s/10)
## Job: test / test / test (ubuntu-24.04-arm)
- Total Job Duration: 2m18s
- The number of Job Executions: 10
- Average Job Duration: 14s
- Slowest jobs: [17s](https://github.com/aquaproj/aqua-registry/actions/runs/18984615607/job/54225335036), [17s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223445803), [16s](https://github.com/aquaproj/aqua-registry/actions/runs/18984012981/job/54223408487)
### Slow steps
1. Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8: 4s (36s/10)
2. Set up job: 3s (27s/10)
3. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 2s (22s/10)
## Job: test / test / test (ubuntu-24.04)
- Total Job Duration: 2m10s
- The number of Job Executions: 10
- Average Job Duration: 13s
- Slowest jobs: [18s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223445796), [17s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227586003), [17s](https://github.com/aquaproj/aqua-registry/actions/runs/18984012981/job/54223408467)
### Slow steps
1. Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8: 4s (37s/10)
2. Set up job: 3s (29s/10)
3. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 2s (20s/10)
## Job: test / ci-info / ci-info
- Total Job Duration: 1m39s
- The number of Job Executions: 10
- Average Job Duration: 10s
- Slowest jobs: [12s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227570936), [12s](https://github.com/aquaproj/aqua-registry/actions/runs/18985291460/job/54227470381), [12s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223426815)
### Slow steps
1. Run suzuki-shunsuke/ci-info-action/store@ceeb10dd50cd632db31e7eccf92cbbb6856f3191: 2s (23s/10)
2. Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342: 2s (21s/10)
## Job: test / check-files / check-files
- Total Job Duration: 47s
- The number of Job Executions: 10
- Average Job Duration: 5s
- Slowest jobs: [7s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227578964), [6s](https://github.com/aquaproj/aqua-registry/actions/runs/18983245993/job/54220752481), [5s](https://github.com/aquaproj/aqua-registry/actions/runs/18984012981/job/54223401351)
The job has no slow steps
## Job: test / path-filter
- Total Job Duration: 38s
- The number of Job Executions: 10
- Average Job Duration: 4s
- Slowest jobs: [5s](https://github.com/aquaproj/aqua-registry/actions/runs/18985324625/job/54227570885), [5s](https://github.com/aquaproj/aqua-registry/actions/runs/18984615607/job/54225313297), [5s](https://github.com/aquaproj/aqua-registry/actions/runs/18984024497/job/54223426756)
The job has no slow steps

</details>

## Why?

[There are some awesome tools](#similar-works), but they can't retrieve step-level data inside composite actions because [the Workflow Jobs APIs](https://docs.github.com/en/rest/actions/workflow-jobs) don’t include those data.
On the other hand, ghaperf can detect bottlenecks within composite actions.
It retrieves job logs via the API, parses them, and extracts data from all log groups — including steps within composite actions.

## Install

[Please see INSTALL.md.](INSTALL.md)

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

2. Run against a workflow run:

A workflow run id is needed.

```sh
run_id=$(gh run list -R suzuki-shunsuke/tfaction -w test.yaml -s completed -L 1 --json databaseId -q ".[0].databaseId")
```

```sh
ghaperf --repo "suzuki-shunsuke/tfaction" --run-id "$run_id" --threshold 1s
```

3. Run against a job:

job id is needed.

```sh
job_id=$(gh api \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "/repos/suzuki-shunsuke/tfaction/actions/runs/${run_id}/jobs?per_page=1" \
  -q ".jobs[0].id")
```

```sh
ghaperf --repo suzuki-shunsuke/tfaction --job-id "$job_id" --threshold 1s
```

4. Run against workflow runs:

Note that this command may take longer for larger count values.

```sh
ghaperf --repo suzuki-shunsuke/ghaperf --workflow test.yaml --count 10 --threshold 1s
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

## Configuration file

Configuration files (`--config <configuration file path>`) are used to filter jobs by job name and normalize matrix job names.
Configuration files and all settings are optional.

```sh
ghaperf --repo suzuki-shunsuke/ghaperf --workflow test.yaml --count 10 --config ghaperf.yaml
```

ghaperf.yaml:

- `job_names`: (Optional) A list of job name glob patterns. Only jobs matching with any patterns are used.
- `excluded_job_names`: (Optional) A list of job name glob patterns. Only jobs not matching with all patterns are used.
- `job_name_mappings`: (Optional) A map of job name glob patterns and normalized job names. This is used to normalize matrix job names.

```yaml
job_names:
  - "test / test / test *"
# excluded_jobs_name:
#   - "test / test / test *"
job_name_mappings:
  "test / test / test *": "test/ test / test"
```

### JSON Schema of Configuration files

[ghaperf.json](json-schema/ghaperf.json)

If you look for a CLI tool to validate configuration with JSON Schema, [ajv-cli](https://ajv.js.org/packages/ajv-cli.html) is useful.

```sh
ajv --spec=draft2020 -s json-schema/ghaperf.json -d pinact.yaml
```

#### Input Complementation by YAML Language Server

[Please see the comment too.](https://github.com/szksh-lab/.github/issues/67#issuecomment-2564960491)

Version: `main`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/ghaperf/main/json-schema/ghaperf.json
```

Or pinning version:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/ghaperf/v0.0.3/json-schema/ghaperf.json
```

## Usage

[Please see here.](USAGE.md)

## Cache

ghaperf caches raw data of completed workflow runs and jobs in the cache directory `${XDG_CACHE_HOME:-${HOME}/.cache}/ghaperf/`.

## :warning: Note

1. ghaperf gets job logs by [GitHub API](https://docs.github.com/en/rest/actions/workflow-jobs#download-job-logs-for-a-workflow-run), but if jobs aren't completed or completed just now, the API would fail.
1. The specification of log format isn't published, so ghaperf may fail to parse logs due to unexpected specification, and it may get broken suddenly due to changes of the log specification.
1. [By default, log files generated by workflows are retained for 90 days before they are automatically deleted.](https://docs.github.com/en/organizations/managing-organization-settings/configuring-the-retention-period-for-github-actions-artifacts-and-logs-in-your-organization)

## Similar Works

- [GitHub Actions Performance Metrics](https://docs.github.com/en/actions/concepts/metrics#about-github-actions-performance-metrics)
- https://github.com/Kesin11/actions-timeline
- https://github.com/Kesin11/github_actions_otel_trace
- https://github.com/Kesin11/CIAnalyzer
- https://github.com/inception-health/otel-export-trace-action
- https://github.com/runforesight/workflow-telemetry-action
- https://github.com/paper2/github-actions-opentelemetry
