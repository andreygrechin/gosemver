package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/andreygrechin/gosemver/internal/commands"
	"github.com/andreygrechin/gosemver/internal/semver"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *semver.SemVer
		wantErr bool
	}{
		// Valid versions according to SemVer 2.0.0
		{"basic version", "1.9.0", &semver.SemVer{Major: 1, Minor: 9, Patch: 0}, false},
		{"with v prefix", "v2.0.0", &semver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"with V prefix", "V2.0.0", &semver.SemVer{Major: 2, Minor: 0, Patch: 0}, false},
		{"with prerelease", "1.0.0-alpha", &semver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha"}, false},
		{"with build", "1.0.0+001", &semver.SemVer{Major: 1, Minor: 0, Patch: 0, Build: "001"}, false},
		{"with prerelease and build", "1.0.0-alpha+001", &semver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha", Build: "001"}, false},
		{"complex prerelease", "1.0.0-alpha.1.beta.11", &semver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.1.beta.11"}, false},
		{"complex build", "1.0.0+20130313144700", &semver.SemVer{Major: 1, Minor: 0, Patch: 0, Build: "20130313144700"}, false},
		{"complex both", "1.0.0-beta.11+exp.sha.5114f85", &semver.SemVer{Major: 1, Minor: 0, Patch: 0, Prerelease: "beta.11", Build: "exp.sha.5114f85"}, false},

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
			got, err := semver.ParseSemVer(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch ||
				got.Prerelease != tt.want.Prerelease || got.Build != tt.want.Build {
				t.Errorf("parseSemVer() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCompareSemVer(t *testing.T) {
	tests := []struct {
		name  string
		v1    string
		v2    string
		want  int
		setup func() (*semver.SemVer, *semver.SemVer)
	}{
		// Major, minor, patch comparisons
		{"major different", "2.0.0", "1.0.0", 1, nil},
		{"minor different", "1.2.0", "1.1.0", 1, nil},
		{"patch different", "1.0.2", "1.0.1", 1, nil},

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
			got, err := commands.CompareSemVer(tt.v1, tt.v2)
			if err != nil {
				t.Errorf("compareSemVer() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("compareSemVer(%s, %s) = %v, want %v", tt.v1, tt.v2, got, tt.want)
			}

			// Test symmetry: if a > b then b < a
			reverse, err := commands.CompareSemVer(tt.v2, tt.v1)
			if err != nil {
				t.Errorf("compareSemVer() error = %v", err)
				return
			}
			if got != 0 && reverse != -got {
				t.Errorf("compareSemVer symmetry failed: %s vs %s: %d and %d", tt.v1, tt.v2, got, reverse)
			}
		})
	}
}

func TestGetSemVer(t *testing.T) {
	tests := []struct {
		name       string
		subCommand string
		version    string
		want       *semver.SemVer
		wantErr    bool
	}{
		// Major component tests
		{"get major basic", "major", "1.2.3", &semver.SemVer{Major: 1}, false},
		{"get major with prerelease", "major", "2.0.0-alpha", &semver.SemVer{Major: 2}, false},
		{"get major with build", "major", "3.0.0+build.123", &semver.SemVer{Major: 3}, false},

		// Minor component tests
		{"get minor basic", "minor", "1.2.3", &semver.SemVer{Minor: 2}, false},
		{"get minor zero", "minor", "1.0.0", &semver.SemVer{Minor: 0}, false},
		{"get minor with metadata", "minor", "1.2.0-beta+build", &semver.SemVer{Minor: 2}, false},

		// Patch component tests
		{"get patch basic", "patch", "1.2.3", &semver.SemVer{Patch: 3}, false},
		{"get patch zero", "patch", "1.2.0", &semver.SemVer{Patch: 0}, false},
		{"get patch complex version", "patch", "1.2.3-alpha.1+build.123", &semver.SemVer{Patch: 3}, false},

		// Prerelease component tests
		{"get prerelease basic", "prerelease", "1.2.3-alpha", &semver.SemVer{Prerelease: "alpha"}, false},
		{"get prerelease empty", "prerelease", "1.2.3", &semver.SemVer{Prerelease: ""}, false},
		{"get prerelease complex", "prerelease", "1.2.3-alpha.1.beta", &semver.SemVer{Prerelease: "alpha.1.beta"}, false},

		// Prerel component tests
		{"get prerelease basic", "prerel", "1.2.3-alpha", &semver.SemVer{Prerelease: "alpha"}, false},
		{"get prerelease empty", "prerel", "1.2.3", &semver.SemVer{Prerelease: ""}, false},
		{"get prerelease complex", "prerel", "1.2.3-alpha.1.beta", &semver.SemVer{Prerelease: "alpha.1.beta"}, false},

		// Build component tests
		{"get build basic", "build", "1.2.3+build", &semver.SemVer{Build: "build"}, false},
		{"get build empty", "build", "1.2.3", &semver.SemVer{Build: ""}, false},
		{"get build with prerelease", "build", "1.2.3-alpha+build.123", &semver.SemVer{Build: "build.123"}, false},

		// Release component tests
		{"get release basic", "release", "1.2.3", &semver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},
		{"get release with prerelease", "release", "1.2.3-alpha", &semver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},
		{"get release with build", "release", "1.2.3+build", &semver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},
		{"get release complex", "release", "1.2.3-alpha.1+build.123", &semver.SemVer{Major: 1, Minor: 2, Patch: 3}, false},

		// Error cases
		{"invalid subcommand", "invalid", "1.2.3", nil, true},
		{"invalid version", "major", "invalid", nil, true},
		{"empty version", "major", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commands.GetSemVer(tt.subCommand, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Compare only the relevant field based on subCommand
			switch tt.subCommand {
			case "major":
				if got.Major != tt.want.Major {
					t.Errorf("GetSemVer() major = %v, want %v", got.Major, tt.want.Major)
				}
			case "minor":
				if got.Minor != tt.want.Minor {
					t.Errorf("GetSemVer() minor = %v, want %v", got.Minor, tt.want.Minor)
				}
			case "patch":
				if got.Patch != tt.want.Patch {
					t.Errorf("GetSemVer() patch = %v, want %v", got.Patch, tt.want.Patch)
				}
			case "prerelease":
				if got.Prerelease != tt.want.Prerelease {
					t.Errorf("GetSemVer() prerelease = %v, want %v", got.Prerelease, tt.want.Prerelease)
				}
			case "build":
				if got.Build != tt.want.Build {
					t.Errorf("GetSemVer() build = %v, want %v", got.Build, tt.want.Build)
				}
			}
		})
	}
}

func TestBumpPrerel(t *testing.T) {
	tests := []struct {
		name     string
		proto    string
		existing string
		want     string
		wantErr  bool
	}{
		// Direct replacement (no trailing dot)
		{"direct replace", "beta", "alpha", "beta", false},
		{"replace empty", "beta", "", "beta", false},

		// Numeric bumping (trailing dot)
		{"init numeric", "alpha.", "", "alpha1", false},
		{"bump existing", "alpha.", "alpha1", "alpha2", false},
		{"change prefix", "beta.", "alpha1", "beta1", false},

		// Auto-bump ("+.")
		{"auto-bump empty", "+.", "", "1", false},
		{"auto-bump numeric", "+.", "1", "2", false},
		{"auto-bump alpha", "+.", "alpha", "alpha1", false},
		{"auto-bump alpha-numeric", "+.", "alpha1", "alpha2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commands.BumpPrerel(tt.proto, tt.existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("bumpPrerel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bumpPrerel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExitCodes(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Run the actual command that might exit
		args := strings.Split(os.Getenv("TEST_ARGS"), " ")
		os.Args = append([]string{"gosemver"}, args...)
		main()
		return
	}

	tests := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{"valid version validate", []string{"validate", "1.0.0"}, 0},
		{"invalid version validate", []string{"validate", "not.a.version"}, 1}, // validate prints "invalid" and exits -1 which becomes 255
		{"valid version compare", []string{"compare", "1.0.0", "2.0.0"}, 0},
		{"invalid version compare", []string{"compare", "not.a.version", "2.0.0"}, 1},
		{"valid version bump", []string{"bump", "major", "1.0.0"}, 0},
		{"invalid version bump", []string{"bump", "major", "not.a.version"}, 1},
		{"valid version get", []string{"get", "major", "1.0.0"}, 0},
		{"invalid version get", []string{"get", "major", "not.a.version"}, 1},
		{"help command", []string{"--help"}, 0},
		{"version command", []string{"--version"}, 0},
		{"unknown command", []string{"unknown"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0])
			cmd.Env = append(os.Environ(),
				"BE_CRASHER=1",
				"TEST_ARGS="+strings.Join(tt.args, " "))
			err := cmd.Run()

			// Get the exit code
			exitCode := 0
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				}
			}

			if exitCode != tt.wantCode {
				t.Errorf("got exit code %d, want %d", exitCode, tt.wantCode)
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
		{"missing parts", "1.0", false},
		{"invalid major", "x.0.0", false},
		{"leading zeros", "01.0.0", false},
		{"invalid prerelease", "1.0.0-alpha@", false},
		{"invalid build", "1.0.0+build@", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semver.SemverRegexp.MatchString(tt.version); got != tt.want {
				t.Errorf("validate(%v) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
