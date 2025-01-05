package cmd

import (
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
		fmt.Printf("error: %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
