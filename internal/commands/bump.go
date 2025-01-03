package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/andreygrechin/gosemver/internal/semver"
)

// bumpSemVer bumps a version with major/minor/patch/prerel/build/release logic.
func BumpSemVer(subCommand, version string) (*semver.SemVer, error) {
	ver, err := semver.ParseSemVer(version)
	if err != nil {
		return nil, err
	}
	switch subCommand {
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
	case "release":
		// remove prerelease & build
		ver.Prerelease = ""
		ver.Build = ""
	// case "prerel":
	// 	newPrere, err := bumpPrerel(subVersion, ver.Prerelease)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ver.Prerelease = newPrere
	// 	ver.Build = ""
	// case "build":
	// 	// sets the build part; leaves the prerelease alone
	// 	ver.Build = subVersion
	default:
		return nil, errors.New("unknown bump command")
	}
	return ver, nil
}

// bumpPrerel replicates the logic of bumping a prerelease based on a "prototype" argument.
// If prototype doesn't end in '.', it simply replaces. If it ends in '.', we bump or initialize
// a numeric suffix. If prototype is "+." (the script's convention), it means there's no user
// prototype, so we just bump or initialize the existing pre-release numeric field.
func BumpPrerel(proto, existing string) (string, error) {
	// If no trailing dot => direct replace
	if !strings.HasSuffix(proto, ".") {
		return proto, nil
	}

	// remove the trailing dot
	proto = strings.TrimSuffix(proto, ".")

	// If user provided no <prerel> argument => sub_version is "+.", meaning we either bump the
	// existing numeric or append "1"
	if proto == "+" {
		// bump the existing numeric (if any)
		return bumpExistingNumeric(existing), nil
	}

	// otherwise, we extract prefix + numeric from existing
	prefix, numeric := extractPrerelParts(existing)
	// If prefix is different from proto => set new prefix + start at 1
	if prefix != proto {
		return fmt.Sprintf("%s1", proto), nil
	}
	// If we already have a numeric => bump it
	if numeric != "" {
		oldNum, _ := strconv.Atoi(numeric)
		oldNum++
		return fmt.Sprintf("%s%d", prefix, oldNum), nil
	}
	// else no numeric => start at 1
	return fmt.Sprintf("%s1", prefix), nil
}

// bumpExistingNumeric tries to bump the numeric suffix in an existing prerelease.
func bumpExistingNumeric(existing string) string {
	prefix, numeric := extractPrerelParts(existing)
	if numeric != "" {
		oldNum, _ := strconv.Atoi(numeric)
		oldNum++
		return fmt.Sprintf("%s%d", prefix, oldNum)
	}
	// if there's no numeric part, append "1"
	return fmt.Sprintf("%s1", prefix)
}

// extractPrerelParts extracts the prefix part (could include letters, dots, hyphens)
// and the trailing numeric part if any.
func extractPrerelParts(prerel string) (string, string) {
	// Regex logic from the script:
	// PREFIX_ALPHANUM='[.0-9A-Za-z-]*[.A-Za-z-]'
	// DIGITS='[0-9][0-9]*'
	// EXTRACT_REGEX="^(${PREFIX_ALPHANUM})*(${DIGITS})$"
	// We can do a simpler approach in Go: loop from the end looking for digits
	idx := -1
	for i := len(prerel) - 1; i >= 0; i-- {
		if prerel[i] < '0' || prerel[i] > '9' {
			// first non-digit from the end => i+1 is the start of numeric
			idx = i
			break
		}
	}
	if idx == len(prerel)-1 {
		// no trailing digits
		return prerel, ""
	}
	return prerel[:idx+1], prerel[idx+1:]
}
