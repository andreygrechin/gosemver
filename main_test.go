package main

import (
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
		{"valid version compare", []string{"compare", "1.0.0", "2.0.0"}, 0},
		{"invalid version compare", []string{"compare", "not.a.version", "2.0.0"}, 1},
		{"valid version bump", []string{"bump", "major", "1.0.0"}, 0},
		{"invalid version bump", []string{"bump", "major", "not.a.version"}, 1},
		{"valid version get", []string{"get", "major", "1.0.0"}, 0},
		{"invalid version get", []string{"get", "major", "not.a.version"}, 1},
		{"help command", []string{"--help"}, 0},
		{"version command", []string{"version"}, 0},
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
