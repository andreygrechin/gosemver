package semver

import (
	"fmt"
	"strconv"

	"regexp"
)

var SemverRegexp = regexp.MustCompile(`^[vV]?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)` +
	`(?:-((?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(?:\.(?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*))?` +
	`(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// SemVer holds the parsed segments of a semantic version.
type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
	Release    string
}

// SemVerToString converts a SemVer object back to a string.
func SemVerToString(ver *SemVer) string {
	s := fmt.Sprintf("%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
	if ver.Prerelease != "" {
		s += "-" + ver.Prerelease
	}
	if ver.Build != "" {
		s += "+" + ver.Build
	}
	return s
}

// ParseSemVer parses a semver string into a SemVer struct.
func ParseSemVer(version string) (*SemVer, error) {
	matches := SemverRegexp.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("version %s does not match the semver scheme 'X.Y.Z(-PRERELEASE)(+BUILD)'", version)
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	prerelease := matches[4]
	build := matches[5]

	return &SemVer{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease, // might be empty
		Build:      build,      // might be empty
		Release:    fmt.Sprintf("%d.%d.%d", major, minor, patch),
	}, nil
}
