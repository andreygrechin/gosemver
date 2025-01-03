package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/internal/commands"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <version1> <version2>",
	Short: "Show differences between versions",
	Long: `Display the component-wise differences between two semantic versions.
Shows which components (major, minor, patch, pre-release, build) differ
between the two versions.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		otherVersion := args[1]
		semVer, err := commands.CommandDiff(version, otherVersion)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("SemVer: %v\n", semVer)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
