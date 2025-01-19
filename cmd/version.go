package cmd

import (
	"fmt"
	"os"

	c "github.com/andreygrechin/gosemver/internal/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version and metadata information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"Version: %s\nBuild Time: %s\nCommit: %s\n",
			c.Version,
			c.BuildTime,
			c.Commit,
		)
		os.Exit(c.ExitOK)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
