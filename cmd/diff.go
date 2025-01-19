//nolint:dupl
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	c "github.com/andreygrechin/gosemver/internal/config"
	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

const (
	diffNArgs      = 2
	diffNArgsStdin = 1
)

var diffCmd = &cobra.Command{
	Use:   "diff <version> <other_version> | diff -",
	Short: "Find the most significant different version identifier",
	Long: `Find the most significant different between <version> and <other_version> identifiers,
output the identifier to stdout.

The versions can be provided either as two arguments or via stdin when using '-' as the argument. In that case,
versions should be separated by a space. Only one input method can be used at a time.

Examples:
  gosemver diff v0.1.2 v0.2.2
  gosemver diff v0.1.2 v0.1.2-beta1
`,
	Args: cobra.RangeArgs(diffNArgsStdin, diffNArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			diffResult    string
			diffResultErr error
		)
		if len(args) == diffNArgs {
			version := args[0]
			otherVersion := args[1]
			diffResult, diffResultErr = gosemver.CommandDiff(version, otherVersion)
		}

		if len(args) == diffNArgsStdin {
			versions, err := gosemver.GetLastArg(*cmd, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get arguments: %v\n", err)
				os.Exit(c.ExitOtherErrors)
			}
			if versions == "" {
				fmt.Fprintln(os.Stderr, "Error: versions string is empty")
				os.Exit(c.ExitOtherErrors)
			}
			versionsSlice := strings.Split(versions, " ")
			if len(versionsSlice) != diffNArgs {
				fmt.Fprintln(os.Stderr, "Error: two versions should be provided")
				os.Exit(c.ExitOtherErrors)
			}

			diffResult, diffResultErr = gosemver.CommandDiff(versionsSlice[0], versionsSlice[1])
		}
		if diffResultErr != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", diffResultErr)
			if errors.Is(diffResultErr, gosemver.ErrInvalidVersion) {
				os.Exit(c.ExitInvalidSemver)
			}
			os.Exit(c.ExitOtherErrors)
		}

		fmt.Println(diffResult)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
