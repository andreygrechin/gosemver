# gosemver

A command-line tool for semantic version manipulation that follows
[Semantic Versioning 2.0.0](https://semver.org) specification.

## Features

- Validate semantic versions
- Compare two versions
- Find differences between versions
- Extract version identifiers
- Bump version identifiers (major, minor, patch, prerelease)
- JSON output support

## Installation

```bash
go install github.com/andreygrechin/gosemver@latest
```

## Usage

### Validate a Version

```bash
gosemver validate 1.2.3
gosemver validate v1.2.3-beta.1+build.123
```

### Compare Versions

Compare two versions, outputs:
    - `-1` if second version is higher
    - `0` if equal
    - `1` if second version is lower

```bash
gosemver compare v1.0.0 v1.1.0
gosemver compare 1.0.0-alpha 1.0.0
```

### Find Version Differences

Identify the most significant difference between versions:

```bash
gosemver diff v1.0.0 v1.1.0
gosemver diff v1.0.0 v1.0.0-beta1
```

### Get Version Identifiers

Extract specific version identifiers:

```bash
gosemver get major 1.2.3
gosemver get minor 1.2.3
gosemver get patch 1.2.3
gosemver get prerelease 1.2.3-beta.1
gosemver get build 1.2.3+build.123
gosemver get release 1.2.3-beta.1+build.123
gosemver get json 1.2.3-beta.1+build.123
```

### Bump Version Identifiers

Increment version identifiers:

```bash
gosemver bump major 1.2.3
gosemver bump minor 1.2.3
gosemver bump patch 1.2.3
gosemver bump prerelease 1.2.3 --prerelease-id beta
gosemver bump release 1.2.3-beta.1
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

This tool is inspired by [semver-tool](https://github.com/fsaintjacques/semver-tool) by François
Saint-Jacques.
