package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gosemver",
	Version: "0.1.0",
	Short:   "A semantic versioning tool",
	Long: `A Golang conversion of the semver-tool Bash script.
Usage, commands, and behaviors are meant to closely mirror the original script.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
