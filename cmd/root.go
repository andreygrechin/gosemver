package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gosemver",
	Long: `gosemver: Validate, compare, diff, and extract identifiers of semantic versions.

See also:
  - https://semver.org
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		var exitError *exec.ExitError

		fmt.Printf("error: %v\n", err)

		if errors.As(err, &exitError) {
			os.Exit(exitError.ExitCode())
		}

		os.Exit(1)
	}
}
