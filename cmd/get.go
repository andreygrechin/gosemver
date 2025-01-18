package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/config"
	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <semver_id> <version>",
	Short: "Extract a value of a version identifier",
	Long: `Extract a value of a version identifier from <version>, where <semver_id> is ( major | minor | patch |
prerelease | build | release ). Additionally you may use 'json' as <semver_id> to get the whole version as JSON
object.

Examples:
  gosemver get major 0.1.2
  gosemver get prerelease 2.0.0-beta1
`,
	Args: cobra.ExactArgs(2), //nolint:mnd
	Run: func(cmd *cobra.Command, args []string) {
		semverID := args[0]
		version := args[1]
		fullSemver, err := gosemver.GetSemVer(semverID, version)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			if errors.Is(err, gosemver.ErrInvalidVersion) {
				os.Exit(config.ExitInvalidSemver)
			} else {
				os.Exit(config.ExitOtherErrors)
			}

		}
		fmt.Println(fullSemver)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
