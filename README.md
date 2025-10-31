# ghaperf

**ghaperf** analyzes the performance of GitHub Actions workflows using GitHub API and raw job logs. Unlike other tools, it can detect bottlenecks **inside composite actions** by parsing job logs and extracting step-level timing data.

## Why ghaperf?

[Existing tools](#related-projects) rely on the [Workflow Jobs API](https://docs.github.com/en/rest/actions/workflow-jobs), which doesn't include step-level data from composite actions. **ghaperf** solves this by:

- Retrieving and parsing raw job logs via the API
- Extracting timing data from all log groups, including steps within composite actions

## Quick Start

1. [Install ghaperf](INSTALL.md)
1. [(Optional) Set your GitHub Access token to avoid GitHub API rate limit and to access private repositories](#github-access-token)
1. Run ghaperf

```sh
export GITHUB_TOKEN=xxx

# Analyze a workflow across multiple runs
ghaperf --repo aquaproj/aqua-registry --workflow test.yaml --count 10 --threshold 2s
```

<details>
<summary>Example Report Output</summary>

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

## Key Features

- **Deep visibility into composite actions** - Detect bottlenecks inside composite actions that other tools miss
- **Multiple analysis modes** - Analyze workflows, workflow runs, or individual jobs
- **Markdown reports** - Generate shareable performance reports
- **No infrastructure needed** - Just a CLI tool, no backend or metrics storage required
- **Intelligent caching** - Cache GitHub API responses for completed runs to speed up analysis
- **Flexible filtering** - Filter and normalize job names using configuration files

## Installation

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

## Usage

For complete usage documentation, see [USAGE.md](USAGE.md).

### Prerequisites

A GitHub access token is required to avoid API rate limits and access private repositories. See [GitHub Access Token](#github-access-token) for details.

```sh
export GITHUB_TOKEN=xxx
```

### Analysis Modes

**1. Analyze multiple workflow runs** (recommended for performance insights)

```sh
ghaperf \
  --repo suzuki-shunsuke/ghaperf \
  --workflow test.yaml \
  --count 10 \
  --threshold 2s
```

> [!NOTE]
> Higher `--count` values provide better insights but take longer to process.

**2. Analyze a single workflow run**

```sh
ghaperf \
  --repo "suzuki-shunsuke/tfaction" \
  --run-id "<workflow run id>"
```

**3. Analyze a specific job**

```sh
ghaperf \
  --repo suzuki-shunsuke/tfaction \
  --job-id "<workflow job id>"
```

**4. Analyze a local log file**

```sh
ghaperf --log-file path/to/job.log
```

See the [example log file](testdata/log.txt) for format reference.

## Configuration

### Environment Variables

 Variable | Description | Default
---|---|---
`GHAPERF_LOG_LEVEL` | Log level: `debug`, `info`, `warn`, `error` | `info`
`GHAPERF_GITHUB_TOKEN` | GitHub access token | -
`GITHUB_TOKEN` | GitHub access token (alternative) | -
`GHAPERF_GHTKN` | Enable [ghtkn](https://github.com/suzuki-shunsuke/ghtkn) integration | `false`
`GHAPERF_THRESHOLD` | Default threshold for slow steps/log groups | `30s`

### GitHub Access Token

A GitHub access token is required to fetch workflow runs, jobs, and logs via the GitHub API.
An access token isn't required for public repositories, but it's recommended to avoid API rate limit.

**Required Permissions:**
- Public repositories: No specific permissions needed
- Private repositories: `Actions: Read` permission

**Setup:**

```sh
# Option 1: Use GITHUB_TOKEN
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxx

# Option 2: Use GHAPERF_GITHUB_TOKEN
export GHAPERF_GITHUB_TOKEN=ghp_xxxxxxxxxxxxx

# Option 3: Use ghtkn integration
export GHAPERF_GHTKN=true
```

### Threshold

ghaperf reports steps and log groups that exceed the specified threshold:

**Threshold Configuration:**
- Default: `30s`
- Set via `--threshold` flag or `GHAPERF_THRESHOLD` environment variable
- Format: [Go duration](https://pkg.go.dev/time#ParseDuration) (e.g., `1s`, `2m30s`)

### Configuration File

Use configuration files to filter jobs by job name and normalize job names.
All settings are optional.

```sh
ghperf --config <configuration file path> ...
```

e.g.

```yaml
# Only analyze jobs matching these glob patterns
job_names:
  - "test / test / test *"

# Exclude jobs matching these glob patterns
# excluded_job_names:
#   - "test / test / test *"

# Normalize matrix job names for aggregation
job_name_mappings:
  "test / test / test *": "test / test / test"
```

**Available Fields:**
- `job_names`: List of glob patterns - only matching jobs are analyzed
- `excluded_job_names`: List of glob patterns - matching jobs are excluded
- `job_name_mappings`: Map of glob patterns to normalized names for matrix jobs

**JSON Schema and Validation:**

The configuration schema is available at [json-schema/ghaperf.json](json-schema/ghaperf.json).

Validate your configuration with [ajv-cli](https://ajv.js.org/packages/ajv-cli.html):

```sh
ajv --spec=draft2020 -s json-schema/ghaperf.json -d ghaperf.yaml
```

**IDE Support:**

Enable auto-completion in your editor by adding this to your config file:

Latest version:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/ghaperf/main/json-schema/ghaperf.json
```

Or pin to a specific version:

```
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/ghaperf/v0.0.3/json-schema/ghaperf.json
```

## Advanced Topics

### Caching

ghaperf automatically caches API responses for completed workflow runs and jobs to improve performance on subsequent analyses.

**Cache location:** `${XDG_CACHE_HOME:-${HOME}/.cache}/ghaperf/`

This speeds up repeated analyses and reduces API calls.

## Important Notes

1. **Log availability timing:** Job logs must be fully processed by GitHub. If a job just completed, the API may not have logs ready yet. Wait a few moments and retry.

2. **Log format changes:** GitHub's log format is not officially documented. ghaperf parses logs based on observed patterns, which may break if GitHub changes the format unexpectedly.

3. **Log retention:** [GitHub retains workflow logs for 90 days by default](https://docs.github.com/en/organizations/managing-organization-settings/configuring-the-retention-period-for-github-actions-artifacts-and-logs-in-your-organization). Analysis of older runs may fail if logs have been deleted.

## Related Projects

While these tools are excellent for analyzing GitHub Actions performance, they don't provide step-level visibility into composite actions like ghaperf does.
We don't aim to replace these tools with ghaperf.
For this reason, it is not our intention to re-implement features already present in other tools, such as CIAnalyzer, in ghaperf.
Rather, ghaperf preserves features that other tools are missing.
Tools like CIAnalyzer are more suitable if you want to know how performance changes over the medium to long term.
On the other hand, ghaperf is more suitable if you want to investigate current performance bottlenecks in more detail.

- [GitHub Actions Performance Metrics](https://docs.github.com/en/actions/concepts/metrics#about-github-actions-performance-metrics) - Official GitHub metrics
- [actions-timeline](https://github.com/Kesin11/actions-timeline) - GitHub Actions to visualize timeline of a workflow job in a job summary
- [CIAnalyzer](https://github.com/Kesin11/CIAnalyzer) - Collect metrics to BigQuery
- Collect metrics to OpenTelemetry
  - [github_actions_otel_trace](https://github.com/Kesin11/github_actions_otel_trace) - Export traces to OpenTelemetry
  - [otel-export-trace-action](https://github.com/inception-health/otel-export-trace-action) - OpenTelemetry trace export
  - [workflow-telemetry-action](https://github.com/runforesight/workflow-telemetry-action)
  - [github-actions-opentelemetry](https://github.com/paper2/github-actions-opentelemetry)
