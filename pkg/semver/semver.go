package semver

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var SemverRegexp = regexp.MustCompile(`^[vV]?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)` +
	`(?:-((?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(?:\.(?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*))?` +
	`(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// SemVer holds the parsed segments of a semantic version.
type SemVer struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	Prerelease string `json:"prerelease"`
	Build      string `json:"build"`
	Release    string `json:"release"`
}

// String converts a SemVer object to a string.
func (v SemVer) String() string {
	return SemVerToString(&v)
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

// IsSemVer checks if a string is a valid semantic version.
func IsSemVer(version string) bool {
	matches := SemverRegexp.FindStringSubmatch(version)
	return matches != nil
}

// CompareSemVer compares two SemVer (ignoring build).
// Returns -1 if left < right, 0 if equal, 1 if left > right.
func CompareSemVer(version, otherVersion string) (int, error) {
	left, err := ParseSemVer(version)
	if err != nil {
		return -1, err
	}
	right, err := ParseSemVer(otherVersion)
	if err != nil {
		return 1, err
	}
	// Compare major, minor, patch
	if left.Major < right.Major {
		return -1, nil
	} else if left.Major > right.Major {
		return 1, nil
	}
	if left.Minor < right.Minor {
		return -1, nil
	} else if left.Minor > right.Minor {
		return 1, nil
	}
	if left.Patch < right.Patch {
		return -1, nil
	} else if left.Patch > right.Patch {
		return 1, nil
	}

	// Compare pre-release
	// If both empty, they are equal
	if left.Prerelease == "" && right.Prerelease == "" {
		return 0, nil
	}
	// If only one is empty, that one is greater (i.e. a version without prerelease is newer)
	if left.Prerelease == "" && right.Prerelease != "" {
		return 1, nil
	}
	if left.Prerelease != "" && right.Prerelease == "" {
		return -1, nil
	}

	// Both are non-empty, compare using semver pre-release rules
	leftFields := strings.Split(left.Prerelease, ".")
	rightFields := strings.Split(right.Prerelease, ".")

	for i := 0; i < len(leftFields) || i < len(rightFields); i++ {
		if i >= len(leftFields) {
			return -1, nil // left is shorter => less
		}
		if i >= len(rightFields) {
			return 1, nil // right is shorter => less
		}
		lf, rf := leftFields[i], rightFields[i]
		// Check if both are numeric
		lNum, lErr := strconv.Atoi(lf)
		rNum, rErr := strconv.Atoi(rf)

		if lErr == nil && rErr == nil {
			// Compare numeric
			if lNum < rNum {
				return -1, nil
			} else if lNum > rNum {
				return 1, nil
			}
			// else equal, keep going
		} else if lErr == nil && rErr != nil {
			// numeric vs string => numeric < string
			return -1, nil
		} else if lErr != nil && rErr == nil {
			// string vs numeric => string > numeric
			return 1, nil
		} else {
			// both string
			if lf < rf {
				return -1, nil
			} else if lf > rf {
				return 1, nil
			}
		}
	}
	return 0, nil
}

// BumpSemVer bumps a version with major/minor/patch/prerelease/build/release logic.
func BumpSemVer(semverId, version, newPrereleaseId string) (*SemVer, error) {
	ver, err := ParseSemVer(version)
	if err != nil {
		return nil, err
	}
	switch semverId {
	case "major":
		ver.Major++
		ver.Minor = 0
		ver.Patch = 0
		ver.Prerelease = ""
		ver.Build = ""
	case "minor":
		ver.Minor++
		ver.Patch = 0
		ver.Prerelease = ""
		ver.Build = ""
	case "patch":
		ver.Patch++
		ver.Prerelease = ""
		ver.Build = ""
	case "prerelease":
		prereleaseId, err := BumpPrerelease(newPrereleaseId, ver.Prerelease)
		if err != nil {
			return nil, err
		}
		ver.Prerelease = prereleaseId
		ver.Build = ""
	case "release":
		ver.Prerelease = ""
		ver.Build = ""
	default:
		return nil, errors.New("unknown bump command")
	}
	return ver, nil
}

// BumpPrerelease replicates the logic of bumping a prerelease based on a "prototype" argument.
// If prototype doesn't end in '.', it simply replaces. If it ends in '.', we bump or initialize
// a numeric suffix. If prototype is "+." (the script's convention), it means there's no user
// prototype, so we just bump or initialize the existing pre-release numeric field.
func BumpPrerelease(newPrereleaseId, currentPrereleaseId string) (string, error) {
	// If user provided a new prerelease ID => use it as is
	if newPrereleaseId != "" {
		return newPrereleaseId, nil
	}

	// If no current prerelease => start with plain 1
	if currentPrereleaseId == "" {
		return "1", nil
	}

	// extract prefix + numeric from existing prerelease and bump it
	prefix, numeric := extractPrereleaseParts(currentPrereleaseId)

	// If we already have a numeric => bump it
	if numeric != "" {
		oldNum, _ := strconv.Atoi(numeric)
		oldNum++
		return fmt.Sprintf("%s%d", prefix, oldNum), nil
	}

	// else no numeric => start at 1
	return fmt.Sprintf("%s1", prefix), nil
}

// CommandDiff returns the difference between two versions (major, minor, patch, prerelease, build).
// If no difference, prints nothing.
func CommandDiff(version, otherVersion string) (string, error) {
	v1, err := ParseSemVer(version)
	if err != nil {
		return "", err
	}
	v2, err := ParseSemVer(otherVersion)
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

// GetSemVer returns the requested SemVer identifier of a version.
// If the SemVer identifier is not found, returns an error.
func GetSemVer(semverId, version string) (string, error) {
	ver, err := ParseSemVer(version)
	if err != nil {
		return "", err
	}
	switch semverId {
	case "major":
		return strconv.Itoa(ver.Major), nil
	case "minor":
		return strconv.Itoa(ver.Minor), nil
	case "patch":
		return strconv.Itoa(ver.Patch), nil
	case "prerelease":
		return ver.Prerelease, nil
	case "release":
		return ver.Release, nil
	case "build":
		return ver.Build, nil
	case "json":
		jsonBytes, err := json.Marshal(ver)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	default:
		return "", errors.New("unknown get command")
	}
}
