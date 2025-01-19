package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	c "github.com/andreygrechin/gosemver/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gosemver",
	Long: `gosemver: A command-line utility and a library for validating, comparing, and manipulating semantic
versions, fully adhering to the Semantic Versioning 2.0.0 specification.

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

		os.Exit(c.ExitOtherErrors)
	}
}
