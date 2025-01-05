package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/pkg/semver"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <version>",
	Short: "Validate a semantic version",
	Long: `Validate if a provided semantic version string complies with semver 2.0.0 specification. Exits with
status 0 if valid, 1 if invalid. Prints "valid" or "invalid" to stdout.

Examples:
  gosemver validate 1.2.3
  gosemver validate v1.2.3-beta.1+build.123
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		if semver.IsSemVer(version) {
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
