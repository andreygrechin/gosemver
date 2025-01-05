package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/pkg/semver"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <version> <other_version>",
	Short: "Find the most significant different version identifier",
	Long: `Find the most significant different version identifier of <version> and <other_version>,
output the identifier to stdout.

Examples:
  gosemver diff v0.1.2 v0.2.2
  gosemver diff v0.1.2 v0.1.2-beta1
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		otherVersion := args[1]
		semVer, err := semver.CommandDiff(version, otherVersion)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(semVer)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
