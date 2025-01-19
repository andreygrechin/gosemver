package cmd

import (
	"fmt"
	"os"

	c "github.com/andreygrechin/gosemver/internal/config"
	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <version|->",
	Short: "Validate a semantic version",
	Long: `Validate whether a provided version string complies with the Semantic Versioning 2.0.0 specification.
Exits with status 0 if valid, 1 if invalid. Prints "valid" or "invalid" to stdout.

The version can be provided either as an argument or via stdin when using '-' as the argument.
Only one input method can be used at a time.

Examples:
  gosemver validate 1.2.3
  gosemver validate v1.2.3-beta.1+build.123
  echo "1.2.3" | gosemver validate -
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var version string
		version, err := gosemver.GetLastArg(*cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get arguments: %v\n", err)
			os.Exit(c.ExitOtherErrors)
		}

		if version == "" {
			fmt.Fprintln(os.Stderr, "Error: version string is empty")
			os.Exit(c.ExitOtherErrors)
		}

		if gosemver.IsSemVer(version) {
			fmt.Println("valid")
			os.Exit(c.ExitOK)
		}

		fmt.Fprintln(os.Stderr, "invalid")
		os.Exit(c.ExitInvalidSemver)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
