package commands

import "github.com/andreygrechin/gosemver/internal/semver"

// commandDiff returns the difference between two versions (major, minor, patch, prerelease, build).
// If no difference, prints nothing.
func CommandDiff(version, otherVersion string) (string, error) {
	v1, err := semver.ParseSemVer(version)
	if err != nil {
		return "", err
	}
	v2, err := semver.ParseSemVer(otherVersion)
	if err != nil {
		return "", err
	}
	if v1.Major != v2.Major {
		return "major", nil
	} else if v1.Minor != v2.Minor {
		return "minor", nil
	} else if v1.Patch != v2.Patch {
		return "patch", nil
	} else if v1.Prerelease != v2.Prerelease {
		return "prerelease", nil
	} else if v1.Build != v2.Build {
		return "build", nil
	}
	return "equal", nil
}
