package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

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
		{"invalid version validate", []string{"validate", "not.a.version"}, 1},
		{"invalid version validate", []string{"validate", ""}, 2},

		{"valid version compare", []string{"compare", "1.0.0", "2.0.0"}, 0},
		{"invalid version compare", []string{"compare", "not.a.version", "2.0.0"}, 1},
		{"invalid version compare", []string{"compare", "1.2.3"}, 2},
		{"invalid version compare", []string{"compare", "1.2.3 1.2.4"}, 0},
		{"invalid version compare", []string{"compare", "1.2.3 1.2.4", "-"}, 2},
		{"invalid version compare", []string{"compare", ""}, 2},

		{"valid version bump", []string{"bump", "major", "1.0.0"}, 0},
		{"invalid version bump", []string{"bump", "major", "not.a.version"}, 1},
		{"invalid version bump with prerelease flag", []string{"bump", "major", "--prerelease", "beta"}, 2},
		{"valid version bump with prerelease flag", []string{"bump", "prerelease", "--prerelease", "beta", "1.2.3"}, 0},
		{"invalid version bump with prerelease flag", []string{"bump", "prerelease", "--prerelease", "be++ta", "1.2.3"}, 1},
		{"invalid version bump1", []string{"bump", "major", "-"}, 2},
		{"invalid version bump2", []string{"bump", "major", ""}, 2},
		{"invalid version bump2", []string{"bump", "major", "--prerelease=beta", "1.2.3"}, 2},
		{"invalid version bump2", []string{"bump", "major", "--build=build123", "1.2.3"}, 2},

		{"valid version get", []string{"get", "major", "1.0.0"}, 0},
		{"invalid version get", []string{"get", "major", "not.a.version"}, 1},
		{"valid version diff", []string{"diff", "v1.2.3", "1.2.4"}, 0},
		{"invalid version diff", []string{"diff", "v1.2.3", "01"}, 1},
		{"invalid version diff", []string{"diff", "1.2.3"}, 2},
		{"invalid version diff", []string{"diff", "1.2.3 1.2.4"}, 0},
		{"invalid version diff", []string{"diff", "1.2.3 1.2.4", "-"}, 2},
		{"invalid version diff", []string{"diff", ""}, 2},

		{"help command", []string{"--help"}, 0},

		{"version command", []string{"version"}, 0},

		{"unknown command", []string{"unknown"}, 2},
	}

	binaryPath, err := os.Executable()
	if err != nil {
		t.Fatalf("failed to get executable path: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath)

			cmd.Env = append(os.Environ(),
				"BE_CRASHER=1",
				"TEST_ARGS="+strings.Join(tt.args, " "))
			err := cmd.Run()

			// Get the exit code
			exitCode := 0

			if err != nil {
				var exitError *exec.ExitError

				if errors.As(err, &exitError) {
					exitCode = exitError.ExitCode()
				}
			}

			if exitCode != tt.wantCode {
				t.Errorf("got exit code %d, want %d", exitCode, tt.wantCode)
			}
		})
	}
}
