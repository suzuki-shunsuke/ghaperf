# Install

ghaperf is written in Go. So you only have to install a binary in your `PATH`.

## Supported OS and architecture

ghaperf supports the following OS and architecture officially.

- linux amd64 and arm64
- darwin (macOS) arm64

[You can build ghaperf yourself from the source code for other platforms.](#build-an-executable-binary-from-source-code-yourself-using-go)

## How To Install

There are some ways to install ghaperf.

1. [Homebrew](#homebrew)
1. [aqua](#aqua)
1. [GitHub Releases](#github-releases)
1. [Build an executable binary from source code yourself using Go](#build-an-executable-binary-from-source-code-yourself-using-go)

## Homebrew

You can install ghaperf using [Homebrew](https://brew.sh/).

```sh
brew install suzuki-shunsuke/ghaperf/ghaperf --cask
```

## aqua

[aqua-registry >= v4.431.0 is required](https://github.com/aquaproj/aqua-registry/releases/tag/v4.431.0).
You can install ghaperf using [aqua](https://aquaproj.github.io/).

```sh
aqua g -i suzuki-shunsuke/ghaperf
```

## Build an executable binary from source code yourself using Go

```sh
go install github.com/suzuki-shunsuke/ghaperf/cmd/ghaperf@latest
```

## GitHub Releases

You can download an asset from [GitHub Releases](https://github.com/suzuki-shunsuke/ghaperf/releases).
Please unarchive it and install a pre built binary into `$PATH`. 

### Verify downloaded assets from GitHub Releases

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [Cosign](https://github.com/sigstore/cosign)
1. [Packslip](https://packslip.dev/) `ghaperf >= v1.0.3`
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier) `ghaperf <= v1.0.1`

### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v1.0.3
asset=ghaperf_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/ghaperf "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/ghaperf \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

### 2. Cosign

You can install Cosign by aqua.
Cosign v2.4.2 or later is required.

```sh
aqua g -i sigstore/cosign
```

From ghaperf v1.0.0, the bundle file is `ghaperf_checksums.txt.sigstore.json`.

```sh
version=v1.0.3
checksum_file="ghaperf_checksums.txt"
asset=ghaperf_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/ghaperf \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.sigstore.json"
cosign verify-blob \
  "$checksum_file" \
  --bundle "${checksum_file}.sigstore.json" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

For ghaperf v0.1.1 or older, the bundle file is `ghaperf_checksums.txt.bundle`.

```sh
version=v0.1.1
checksum_file="ghaperf_checksums.txt"
asset=ghaperf_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/ghaperf \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.bundle"
cosign verify-blob \
  "$checksum_file" \
  --bundle "${checksum_file}.bundle" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

### 3. Packslip

`ghaperf >= v1.0.3`

You can install packslip by aqua.

```sh
aqua g -i jdx/packslip
```

```sh
version=v1.0.3
asset=ghaperf_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/ghaperf "$version" -p "$asset" -p packslip.sigstore.json
packslip verify packslip.sigstore.json \
    --pin ps1_5blxwuxlg3dkisoamdm3w6gl6q \
    --identity-prefix https://github.com/suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@ \
    --artifact "$asset"
```

- `--pin` requires the release to be signed from the repository `suzuki-shunsuke/ghaperf`, identified by its repository ID rather than its name
- `--identity-prefix` requires the release to be signed by the release workflow of [suzuki-shunsuke/go-release-workflow](https://github.com/suzuki-shunsuke/go-release-workflow)

On success, packslip prints a line starting with `ok:`.
Please make sure it shows the project and version you meant, e.g. `github.com/suzuki-shunsuke/ghaperf 1.0.3`.

### 4. slsa-verifier

> [!WARNING]
> As of ghaperf v1.0.2, `multiple.intoto.jsonl` hasn't been published anymore.

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v1.0.1
asset=ghaperf_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/ghaperf "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/ghaperf \
  --source-tag "$version"
```
