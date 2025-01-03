package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/commands"
	"github.com/andreygrechin/gosemver/internal/semver"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <component> <version>",
	Short: "Extract a version component",
	Long: `Extract and display a specific component from a semantic version.
Valid components are:
  major       - Major version number
  minor       - Minor version number
  patch       - Patch version number
  prerel      - Pre-release version
  prerelease  - Alias for prerel
  build       - Build metadata
  release     - Version without pre-release or build metadata`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		subCommand := args[0]
		version := args[1]
		semVer, err := commands.GetSemVer(subCommand, version)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(semver.SemVerToString(semVer))

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
