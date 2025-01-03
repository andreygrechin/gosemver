package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/semver"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <version>",
	Short: "Check if version is valid",
	Long: `Validate if a version string complies with semver 2.0.0 specification.
Exits with status 0 if valid, 1 if invalid.
Prints "valid" or "invalid" to stdout.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		if semver.SemverRegexp.MatchString(version) {
			fmt.Println("valid")
		} else {
			fmt.Println("invalid")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
