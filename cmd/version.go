package cmd

import (
	"github.com/andreygrechin/gosemver/internal/commands"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Long:  "Display version, build time, and commit information for this application.",
	Run: func(cmd *cobra.Command, args []string) {
		commands.Version(args)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
