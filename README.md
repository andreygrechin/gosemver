# gosemver

A command-line utility and a library for validating, comparing, and manipulating semantic versions,
fully adhering to the [Semantic Versioning 2.0.0](https://semver.org) specification. Written in Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/andreygrechin/gosemver)](https://goreportcard.com/report/github.com/andreygrechin/gosemver)
[![codecov](https://codecov.io/gh/andreygrechin/gosemver/graph/badge.svg?token=789FTQPOB0)](https://codecov.io/gh/andreygrechin/gosemver)
[![Go Reference](https://pkg.go.dev/badge/github.com/andreygrechin/gosemver.svg)](https://pkg.go.dev/github.com/andreygrechin/gosemver)
[![license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/andreygrechin/gosemver/blob/main/licenses/LICENSE.MIT)
[![license](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/andreygrechin/gosemver/blob/main/licenses/LICENSE.Apache-2.0)

## Features

- Validate semantic versions
- Compare two versions
- Find differences between versions
- Extract version identifiers
- Bump version identifiers (major, minor, patch, prerelease)
- JSON output support

## Installation

### go install

```shell
go install github.com/andreygrechin/gosemver@latest
```

### Homebrew tap

You may also install the latest version of `gosemver` using the Homebrew tap:

```shell
brew install andreygrechin/tap/gosemver

# to update, run
brew update
brew upgrade gosemver
```

### Containers

```shell
$ docker run --rm ghcr.io/andreygrechin/gosemver:latest validate 1.2.3
valid
```

### Manually

Download the pre-compiled binaries from
[the releases page](https://github.com/andreygrechin/gosemver/releases/) and copy them to a desired
location.

## Usage

### Validate a Version

```shell
$ gosemver validate 1.2.3
valid

$ gosemver validate v1.2.3-beta.1+build.123
valid

$ gosemver validate 1.2
invalid # also returns exit code 1
```

### Compare Versions

Compare two versions, outputs:

- `-1` if the first version is lower
- `0` if equal
- `1` if the first version is higher

```shell
$ gosemver compare v1.0.0 v1.2.3
-1

$ gosemver compare 1.0.0 1.0.0-alpha
1
```

### Find Version Differences

Identify the most significant difference between versions:

```shell
$ gosemver diff v1.0.0 v1.1.0
minor

$ gosemver diff v1.0.0 v1.0.0-beta1
prerelease
```

### Get Version Identifiers

Extract specific version identifiers:

```shell
$ gosemver get major 1.2.3
1

$ gosemver get prerelease 1.2.3-beta.1
beta.1

$ gosemver get build 1.2.3+build.123
build.123

$ gosemver get release 1.2.3-beta.1+build.123
1.2.3

$ gosemver get json 1.2.3-beta.1+build.123
{"major":1,"minor":2,"patch":3,"prerelease":"beta.1","build":"build.123","release":"1.2.3"}
```

### Bump Version Identifiers

Increment version identifiers:

```shell
$ gosemver bump major 1.2.3
2.0.0

$ gosemver bump prerelease 1.2.3-beta1
1.2.3-beta2

$ gosemver bump release 1.2.3-beta.1
1.2.3
```

## License

This project can be licensed under MIT or the Apache 2.0 licenses — see the
[LICENSE.Apache-2.0](licenses/LICENSE.Apache-2.0) and the [MIT License](licenses/LICENSE.MIT) files.
Choose the most appropriate for your use case.

`SPDX-License-Identifier: MIT OR Apache-2.0`

## Acknowledgments

This tool is inspired by [semver-tool](https://github.com/fsaintjacques/semver-tool) created by
François Saint-Jacques.
