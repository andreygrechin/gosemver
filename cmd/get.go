package cmd

import (
	"errors"
	"fmt"
	"os"

	c "github.com/andreygrechin/gosemver/internal/config"
	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <semver_id> <version>",
	Short: "Extract a value of a version identifier",
	Long: `Extract a value of a version identifier from <version>, where <semver_id> is ( major | minor | patch |
prerelease | build | release ). Additionally you may use 'json' as <semver_id> to get the whole version as JSON
object.

The version can be provided either as an argument or via stdin when using '-' as the argument.
Only one input method can be used at a time.

Examples:
  gosemver get major 0.1.2
  gosemver get prerelease 2.0.0-beta1
`,
	Args: cobra.ExactArgs(2), //nolint:mnd
	Run: func(cmd *cobra.Command, args []string) {
		semverID := args[0]
		version, err := gosemver.GetLastArg(*cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get arguments: %v\n", err)
			os.Exit(c.ExitOtherErrors)
		}
		if version == "" {
			fmt.Fprintln(os.Stderr, "Error: version string is empty")
			os.Exit(c.ExitOtherErrors)
		}
		fullSemver, err := gosemver.GetSemVer(semverID, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if errors.Is(err, gosemver.ErrInvalidVersion) {
				os.Exit(c.ExitInvalidSemver)
			}
			os.Exit(c.ExitOtherErrors)
		}
		fmt.Println(fullSemver)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
