package commands

import (
	"strconv"
	"strings"

	"github.com/andreygrechin/gosemver/internal/semver"
)

// compareSemVer compares two SemVer (ignoring build).
// Returns -1 if left < right, 0 if equal, 1 if left > right.
func CompareSemVer(version, otherVersion string) (int, error) {
	left, err := semver.ParseSemVer(version)
	if err != nil {
		return -1, err
	}
	right, err := semver.ParseSemVer(otherVersion)
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
