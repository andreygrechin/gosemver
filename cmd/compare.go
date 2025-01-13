package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare <version> <other_version>",
	Short: "Compare two semantic versions",
	Long: `Compare <version> and <other_version> semantic versions, output to stdout -1 if <other_version> is
higher, 0 if equal, 1 if lower. Build identifiers of versions is always ignored.

Examples:
  gosemver compare v0.1.2 v0.1.2-beta1
`,
	Args: cobra.ExactArgs(2), //nolint:mnd
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		otherVersion := args[1]
		semVer, err := gosemver.CompareSemVer(version, otherVersion)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(semVer)
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
