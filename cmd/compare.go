package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/commands"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare <version1> <version2>",
	Short: "Compare two semantic versions",
	Long: `Compare two semantic versions according to semver 2.0.0 precedence rules.
Returns:
  -1 if version1 < version2
   0 if version1 = version2
   1 if version1 > version2`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		otherVersion := args[1]
		semVer, err := commands.CompareSemVer(version, otherVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("SemVer: %v\n", semVer)
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
