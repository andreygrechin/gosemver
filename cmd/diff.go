package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/config"
	"github.com/andreygrechin/gosemver/pkg/gosemver"
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
	Args: cobra.ExactArgs(2), //nolint:mnd
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		otherVersion := args[1]
		semVer, err := gosemver.CommandDiff(version, otherVersion)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			if errors.Is(err, gosemver.ErrInvalidVersion) {
				os.Exit(config.ExitInvalidSemver)
			} else {
				os.Exit(config.ExitOtherErrors)
			}
		}
		fmt.Println(semVer)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
