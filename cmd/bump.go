package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/commands"
	"github.com/andreygrechin/gosemver/internal/semver"

	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:   "bump <component> <version>",
	Short: "Increment a version component",
	Long: `Increment a specific component of a semantic version.
Valid components are:
  major       - Increment major version (resets minor and patch to 0)
  minor       - Increment minor version (resets patch to 0)
  patch       - Increment patch version
  prerel      - Set or increment pre-release version
  prerelease  - Alias for prerel
  build       - Set build metadata
  release     - Remove pre-release and build metadata`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		subCommand := args[0]
		version := args[1]
		semVer, err := commands.BumpSemVer(subCommand, version)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(semver.SemVerToString(semVer))
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
}
