package gosemver_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/andreygrechin/gosemver/pkg/gosemver"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *gosemver.SemVer
		wantErr bool
	}{
		// Valid versions according to SemVer 2.0.0
		{"basic version", "1.9.0", &gosemver.SemVer{Major: 1, Minor: 9, Patch: 0}, false},
		{"with v prefix", "v2.0.0", &gosemver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"with V prefix", "V2.0.0", &gosemver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"with prerelease", "1.0.0-alpha", &gosemver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha"}, false},
		{"with build", "1.0.0+001", &gosemver.SemVer{Major: 1, Minor: 0, Patch: 0, Build: "001"}, false},
		{"with prerelease and build", "1.0.0-alpha+001", &gosemver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha", Build: "001"}, false},
		{"complex prerelease", "1.0.0-alpha.1.beta.11", &gosemver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.1.beta.11"}, false},
		{"complex build", "1.0.0+20130313144700", &gosemver.SemVer{Major: 1, Minor: 0, Patch: 0, Build: "20130313144700"}, false},
		{"complex both", "1.0.0-beta.11+exp.sha.5114f85", &gosemver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "beta.11", Build: "exp.sha.5114f85"}, false},

		// Invalid versions
		{"empty string", "", nil, true},
		{"missing minor", "1.0", nil, true},
		{"missing patch", "1", nil, true},
		{"invalid major", "x.0.0", nil, true},
		{"invalid minor", "1.x.0", nil, true},
		{"invalid patch", "1.0.x", nil, true},
		{"leading zeros major", "01.0.0", nil, true},
		{"leading zeros minor", "1.01.0", nil, true},
		{"leading zeros patch", "1.0.01", nil, true},
		{"invalid prerelease chars", "1.0.0-alpha@", nil, true},
		{"invalid build chars", "1.0.0+build@", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gosemver.ParseSemVer(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch ||
				got.Prerelease != tt.want.Prerelease || got.Build != tt.want.Build {
				t.Errorf("ParseSemVer() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCompareSemVer(t *testing.T) {
	tests := []struct {
		name    string
		v1      string
		v2      string
		want    int
		wantErr error
	}{
		// Major, minor, patch comparisons
		{"major different", "2.0.0", "1.0.0", 1, nil},
		{"minor different", "1.2.0", "1.1.0", 1, nil},
		{"patch different", "1.0.2", "1.0.1", 1, nil},

		// Invalid versions
		{"invalid version 1", "1.0.0", "invalid", 0, fmt.Errorf("%w: %s", gosemver.ErrInvalidVersion, "invalid")},
		{"invalid version 2", "invalid", "1.0.0", 0, fmt.Errorf("%w: %s", gosemver.ErrInvalidVersion, "invalid")},

		// Pre-release comparisons (spec item 11)
		{"no prerelease > prerelease", "1.0.0", "1.0.0-alpha", 1, nil},
		{"alpha < beta", "1.0.0-alpha", "1.0.0-beta", -1, nil},
		{"numeric comparison", "1.0.0-alpha.1", "1.0.0-alpha.2", -1, nil},
		{"numeric < non-numeric", "1.0.0-2", "1.0.0-alpha", -1, nil},
		{"shorter < longer", "1.0.0-alpha", "1.0.0-alpha.1", -1, nil},

		// Build metadata should be ignored in precedence
		{"ignore build", "1.0.0+build.1", "1.0.0+build.2", 0, nil},
		{"ignore build with prerelease", "1.0.0-alpha+build.1", "1.0.0-alpha+build.2", 0, nil},

		// Equal versions
		{"exactly equal", "1.0.0", "1.0.0", 0, nil},
		{"equal with build", "1.0.0+build.1", "1.0.0+build.2", 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gosemver.CompareSemVer(tt.v1, tt.v2)
			if err != nil {
				if errors.Is(err, gosemver.ErrInvalidVersion) {
					return
				}

				t.Errorf("CompareSemVer(%s, %s) = %v, want %v", tt.v1, tt.v2, got, tt.wantErr)
				fmt.Println(err)
				return
			}

			reverse, err := gosemver.CompareSemVer(tt.v2, tt.v1)
			if err != nil {
				t.Errorf("CompareSemVer() error = %v", err)
				return
			}

			if got != 0 && reverse != -got {
				t.Errorf("CompareSemVer symmetry failed: %s vs %s: %d and %d", tt.v1, tt.v2, got, reverse)
			}
		})
	}
}

func TestGetSemVer(t *testing.T) {
	tests := []struct {
		name     string
		semverID string
		version  string
		want     string
		wantErr  bool
	}{
		// Major component tests
		{"get major basic", "major", "1.2.3", "1", false},
		{"get major with prerelease", "major", "2.0.0-alpha", "2", false},
		{"get major with build", "major", "3.0.0+build.123", "3", false},

		// Minor component tests
		{"get minor basic", "minor", "1.2.3", "2", false},
		{"get minor zero", "minor", "1.0.0", "0", false},
		{"get minor with metadata", "minor", "1.2.0-beta+build", "2", false},

		// Patch component tests
		{"get patch basic", "patch", "1.2.3", "3", false},
		{"get patch zero", "patch", "1.2.0", "0", false},
		{"get patch complex version", "patch", "1.2.3-alpha.1+build.123", "3", false},

		// Prerelease component tests
		{"get prerelease basic", "prerelease", "1.2.3-alpha", "alpha", false},
		{"get prerelease empty", "prerelease", "1.2.3", "", false},
		{"get prerelease complex", "prerelease", "1.2.3-alpha.1.beta", "alpha.1.beta", false},

		// Build component tests
		{"get build basic", "build", "1.2.3+build", "build", false},
		{"get build empty", "build", "1.2.3", "", false},
		{"get build with prerelease", "build", "1.2.3-alpha+build.123", "build.123", false},

		// Release component tests
		{"get release basic", "release", "1.2.3", "1.2.3", false},
		{"get release with prerelease", "release", "1.2.3-alpha", "1.2.3", false},
		{"get release with build", "release", "1.2.3+build", "1.2.3", false},
		{"get release complex", "release", "1.2.3-alpha.1+build.123", "1.2.3", false},

		// Error cases
		{"invalid semver id", "invalid", "1.2.3", "", true},
		{"invalid version", "major", "invalid", "", true},
		{"empty version", "major", "", "", true},

		// JSON output tests
		{"get json basic", "json", "1.2.3", `{"major":1,"minor":2,"patch":3,"prerelease":"","build":"","release":"1.2.3"}`, false},
		{"get json with prerelease", "json", "1.2.3-alpha", `{"major":1,"minor":2,"patch":3,"prerelease":"alpha","build":"","release":"1.2.3"}`, false},
		{"get json with build", "json", "1.2.3+build", `{"major":1,"minor":2,"patch":3,"prerelease":"","build":"build","release":"1.2.3"}`, false},
		{"get json complex", "json", "1.2.3-alpha+build", `{"major":1,"minor":2,"patch":3,"prerelease":"alpha","build":"build","release":"1.2.3"}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gosemver.GetSemVer(tt.semverID, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("GetSemVer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBumpPrerelease(t *testing.T) {
	tests := []struct {
		name     string
		proto    string
		existing string
		want     string
		wantErr  bool
	}{
		// Direct replacement
		{"direct replace up", "beta", "alpha", "beta", false},
		{"direct replace down", "alpha", "beta", "alpha", false},
		{"replace empty", "beta", "", "beta", false},
		{"auto-bump numeric", "0", "1", "0", false},

		// Incorrect newPrereleaseId
		{"incorrect identifier", "beta.", "", "beta.", false},

		// Auto-bump numeric
		{"auto-bump empty", "", "", "1", false},
		{"auto-bump empty", "", "0", "1", false},
		{"auto-bump empty", "", "1", "2", false},

		// Auto-bump alphanumeric
		{"auto-bump 1", "", "beta0", "beta1", false},
		{"auto-bump 2", "", "beta1", "beta2", false},
		{"auto-bump 3", "", "beta.0", "beta.1", false},
		{"auto-bump 4", "", "beta.1", "beta.2", false},
		{"auto-bump 5", "", "beta.10", "beta.11", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gosemver.BumpNumericSuffix(tt.proto, tt.existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("BumpPrerelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("BumpPrerelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBumpSemVer(t *testing.T) {
	tests := []struct {
		name     string
		semverID string
		version  string
		want     *gosemver.SemVer
		wantErr  bool
	}{
		// Major version bumps
		{"bump major basic", "major", "1.2.3", &gosemver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"bump major with prerelease", "major", "1.2.3-alpha", &gosemver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"bump major with build", "major", "1.2.3+build", &gosemver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"bump major complex", "major", "1.2.3-alpha+build", &gosemver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},

		// Minor version bumps
		{"bump minor basic", "minor", "1.2.3", &gosemver.SemVer{Major: 1, Minor: 3, Patch: 0}, false},
		{"bump minor with prerelease", "minor", "1.2.3-beta", &gosemver.SemVer{Major: 1, Minor: 3, Patch: 0}, false},
		{"bump minor with build", "minor", "1.2.3+build.123", &gosemver.SemVer{Major: 1, Minor: 3, Patch: 0}, false},
		{"bump minor complex", "minor", "1.2.3-beta+build.123", &gosemver.SemVer{Major: 1, Minor: 3, Patch: 0}, false},

		// Patch version bumps
		{"bump patch basic", "patch", "1.2.3", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 4}, false},
		{"bump patch with prerelease", "patch", "1.2.3-rc1", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 4}, false},
		{"bump patch with build", "patch", "1.2.3+sha.xyz", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 4}, false},
		{"bump patch complex", "patch", "1.2.3-rc1+sha.xyz", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 4}, false},

		// Prerelease version (removes build)
		{"bump prerelease no prerelease", "prerelease", "1.2.3", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "1"}, false},
		{"bump prerelease no numeric suffix", "prerelease", "1.2.3-beta", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta1"}, false},
		{"bump prerelease increment numeric suffix", "prerelease", "1.2.3-beta1", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta2"}, false},

		// Build
		{"bump build basic", "build", "1.2.3", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Build: "1"}, false},
		{"bump build basic", "build", "1.2.3+build", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Build: "build1"}, false},
		{"bump build basic", "build", "1.2.3+build1", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Build: "build2"}, false},

		// Release version (removes prerelease and build)
		{"bump release basic", "release", "1.2.3", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},
		{"bump release with prerelease", "release", "1.2.3-alpha", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},
		{"bump release with build", "release", "1.2.3+build", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},
		{"bump release complex", "release", "1.2.3-alpha+build", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},

		// Error cases
		{"invalid command", "invalid", "1.2.3", nil, true},
		{"invalid version", "major", "invalid", nil, true},
		{"empty version", "major", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gosemver.BumpSemVer(tt.semverID, tt.version, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("BumpSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch ||
				got.Prerelease != tt.want.Prerelease || got.Build != tt.want.Build {
				t.Errorf("BumpSemVer() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDiffCommand(t *testing.T) {
	tests := []struct {
		name         string
		version1     string
		version2     string
		expectedDiff string
		wantErr      bool
	}{
		{"equal versions", "10.1.4", "10.1.4", "equal", false},
		{"prerelease difference", "1.0.1-rc1.1.0+build.051", "1.0.1", "prerelease", false},
		{"minor difference with prerelease", "10.1.4-rc4", "10.4.2-rc1", "minor", false},
		{"major difference", "2.0.0", "1.0.0", "major", false},
		{"minor difference", "1.2.0", "1.1.0", "minor", false},
		{"patch difference", "1.0.2", "1.0.1", "patch", false},
		{"invalid version", "not.valid", "1.0.0", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gosemver.CommandDiff(tt.version1, tt.version2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.expectedDiff {
				t.Errorf("CommandDiff() = %v, want %v", got, tt.expectedDiff)
			}
		})
	}
}

func TestIsPrerelease(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		// Valid prereleases
		{"empty string", "", true},
		{"simple alpha", "alpha", true},
		{"alpha with number", "alpha1", true},
		{"number only", "123", true},
		{"complex with dots", "alpha.1.beta.2", true},
		{"with hyphens", "rc-1", true},
		{"mixed case", "Alpha1", true},

		// Invalid prereleases
		{"with spaces", "alpha 1", false},
		{"with special chars", "alpha@1", false},
		{"with plus", "alpha+1", false},
		{"with underscore", "alpha_1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gosemver.IsPrerelease(tt.version); got != tt.want {
				t.Errorf("IsPrerelease(%v) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestIsBuild(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		// Valid builds
		{"empty string", "", true},
		{"simple number", "001", true},
		{"simple alpha", "build", true},
		{"alpha with number", "build123", true},
		{"with dots", "build.123", true},
		{"with hyphens", "build-123", true},
		{"complex", "20130313144700.exp.sha.5114f85", true},

		// Invalid builds
		{"with spaces", "build 123", false},
		{"with special chars", "build@123", false},
		{"with plus", "build+123", false},
		{"with underscore", "build_123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gosemver.IsBuild(tt.version); got != tt.want {
				t.Errorf("IsBuild(%v) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		// Valid versions
		{"basic version", "1.9.0", true},
		{"with v prefix", "v2.0.0", true},
		{"with prerelease", "1.0.0-alpha", true},
		{"with build", "1.0.0+001", true},
		{"with both", "1.0.0-alpha+001", true},
		{"complex prerelease", "1.0.0-alpha.1.beta.11", true},
		{"complex build", "1.0.0+20130313144700", true},
		{"complex both", "1.0.0-beta.11+exp.sha.5114f85", true},

		// Invalid versions
		{"empty string", "", false},
		{"missing semver id", "1.0", false},
		{"invalid major", "x.0.0", false},
		{"leading zeros", "01.0.0", false},
		{"invalid prerelease", "1.0.0-alpha@", false},
		{"invalid build", "1.0.0+build@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gosemver.SemverRegexp.MatchString(tt.version); got != tt.want {
				t.Errorf("validate(%v) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestSemverToString(t *testing.T) {
	tests := []struct {
		name          string
		semverVersion *gosemver.SemVer
		want          string
	}{
		{"basic version", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3}, "1.2.3"},
		{"basic version", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta"}, "1.2.3-beta"},
		{"basic version", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta", Build: "build123"}, "1.2.3-beta+build123"},
		{"basic version", &gosemver.SemVer{Major: 1, Minor: 2, Patch: 3, Build: "build123"}, "1.2.3+build123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gosemver.ToString(tt.semverVersion); got != tt.want {
				t.Errorf("validate(%v) = %v, want %v", tt.semverVersion, got, tt.want)
			}
		})
	}
}

func TestIsSemVer(t *testing.T) {
	tests := []struct {
		name          string
		semverVersion string
		want          bool
	}{
		{"basic version", "1.2.3", true},
		{"basic version", "beta", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gosemver.IsSemVer(tt.semverVersion); got != tt.want {
				t.Errorf("validate(%v) = %v, want %v", tt.semverVersion, got, tt.want)
			}
		})
	}
}
