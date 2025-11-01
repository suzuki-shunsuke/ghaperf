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
<table>
<tr><td>Average Job Duration</td><td>33s (5m26s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987585604/job/54234368246">40s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986937578/job/54232538653">35s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18985653020/job/54228560750">35s</a></td></tr>
</table>

### Slow steps
1. 13s (2m9s/10): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
    1. 8s (1m23s/10): Run $(if($env:AQUA_ROOT_DIR) {echo $env:AQUA_ROOT_DIR} else {echo "$HOME/AppData/Local/aquaproj-aqua/bin"}) | Out-File -FilePath $env:GITHUB_PATH -Encoding utf8 -Append
    2. 5s (42s/9): Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then
2. 7s (1m13s/10): Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8
    1. 2s (16s/7): Setting up auth
    2. 3s (3s/1): Fetching the repository
3. 4s (36s/10): Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8
4. 4s (35s/10): Set up job
## Job: test / test / test (windows-latest, arm64)
<table>
<tr><td>Average Job Duration</td><td>32s (5m24s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233211736">48s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986103791/job/54230012684">39s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986654682/job/54231704451">36s</a></td></tr>
</table>

### Slow steps
1. 14s (2m18s/10): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
    1. 10s (1m30s/9): Run $(if($env:AQUA_ROOT_DIR) {echo $env:AQUA_ROOT_DIR} else {echo "$HOME/AppData/Local/aquaproj-aqua/bin"}) | Out-File -FilePath $env:GITHUB_PATH -Encoding utf8 -Append
    2. 4s (39s/9): Run if [ "${SKIP_INSTALL_AQUA:-}" = true ] && command -v aqua >/dev/null; then
2. 7s (1m11s/10): Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8
    1. 3s (17s/6): Setting up auth
    2. 3s (3s/1): Getting Git version info
    3. 2s (2s/1): Fetching the repository
3. 4s (36s/10): Set up job
4. 3s (30s/10): Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8
## Job: test / test / test (macos-13)
<table>
<tr><td>Average Job Duration</td><td>23s (3m45s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18985653020/job/54228560759">33s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233211747">31s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986937578/job/54232538667">25s</a></td></tr>
</table>

### Slow steps
1. 6s (1m3s/10): Set up job
2. 6s (51s/9): Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8
    1. 5s (31s/6): Run aqua i --test
    2. 2s (14s/6): Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV"
3. 4s (36s/9): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
4. 4s (34s/9): Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8
## Job: test / test / test (macos-14)
<table>
<tr><td>Average Job Duration</td><td>20s (3m17s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986937578/job/54232538680">22s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233211744">21s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986673543/job/54231758517">21s</a></td></tr>
</table>

### Slow steps
1. 6s (58s/10): Set up job
2. 4s (42s/10): Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8
    1. 3s (21s/7): Run aqua i --test
    2. 2s (7s/3): Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV"
3. 3s (32s/10): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
4. 3s (25s/10): Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8
## Job: test / lintnet / lintnet
<table>
<tr><td>Average Job Duration</td><td>15s (2m29s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986937578/job/54232530626">17s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986654682/job/54231696170">17s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986103791/job/54229999783">15s</a></td></tr>
</table>

### Slow steps
1. 6s (56s/10): Run lintnet lint
2. 3s (30s/10): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
3. 2s (20s/10): Post Run actions/checkout@08c6903cd8c0fde910a37f88322edcfb5dd907a8
## Job: test / test / test (ubuntu-24.04-arm)
<table>
<tr><td>Average Job Duration</td><td>15s (2m29s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233211763">18s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986103791/job/54230012676">17s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18985946839/job/54229501665">17s</a></td></tr>
</table>

### Slow steps
1. 3s (32s/10): Set up job
2. 3s (32s/10): Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8
3. 3s (26s/10): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
## Job: test / test / test (ubuntu-24.04)
<table>
<tr><td>Average Job Duration</td><td>14s (2m16s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233211734">18s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986941589/job/54233657194">15s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986654682/job/54231704444">15s</a></td></tr>
</table>

### Slow steps
1. 3s (33s/10): Set up job
2. 3s (32s/10): Run aquaproj/registry-action/test@68f10339de561d67f9acea40b91dc36aa5011ea8
    1. 4s (4s/1): Run aqua i --test
    2. 2s (2s/1): Run aqua exec -- ci-info run | sed -E "s/^export //" >> "$GITHUB_ENV"
## Job: test / ci-info / ci-info
<table>
<tr><td>Average Job Duration</td><td>10s (1m43s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233199054">14s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986103791/job/54229999767">12s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18985946839/job/54229488797">12s</a></td></tr>
</table>

### Slow steps
1. 2s (22s/10): Run suzuki-shunsuke/ci-info-action/store@ceeb10dd50cd632db31e7eccf92cbbb6856f3191
2. 2s (20s/10): Set up job
3. 2s (20s/10): Run aquaproj/aqua-installer@ea518c135a02fc11ff8024364510c181a5c6b342
## Job: test / check-files / check-files
<table>
<tr><td>Average Job Duration</td><td>5s (53s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986673543/job/54231754113">7s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18985946839/job/54229494701">7s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986649585/job/54231685141">6s</a></td></tr>
</table>

The job has no slow steps
## Job: test / path-filter
<table>
<tr><td>Average Job Duration</td><td>4s (40s/10)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987585604/job/54234360739">5s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986673543/job/54231748927">5s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986654682/job/54231696160">5s</a></td></tr>
</table>

The job has no slow steps
## Job: status-check
<table>
<tr><td>Average Job Duration</td><td>2s (2s/1)</td></tr>
<tr><td>Slowest Jobs</td><td><a href="https://github.com/aquaproj/aqua-registry/actions/runs/18986941589/job/54233677718">2s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987585604/job/54234395291">0s</a>, <a href="https://github.com/aquaproj/aqua-registry/actions/runs/18987158767/job/54233253944">0s</a></td></tr>
</table>

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
