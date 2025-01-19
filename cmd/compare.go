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
	compareNArgs      = 2
	compareNArgsStdin = 1
)

var compareCmd = &cobra.Command{
	Use:   "compare <version> <other_version> | compare -",
	Short: "Compare two semantic versions",
	Long: `Compare <version> and <other_version> semantic versions, output to stdout -1 if <other_version> is
higher, 0 if equal, 1 if lower. Build identifiers of versions is always ignored.

The versions can be provided either as two arguments or via stdin when using '-' as the argument. In that case,
versions should be separated by a space. Only one input method can be used at a time.

Examples:
  gosemver compare v0.1.2 v0.1.2-beta1
  gosemver compare v0.1.2 v0.1.2+build1
`,
	Args: cobra.RangeArgs(compareNArgsStdin, compareNArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			compareResult    int
			compareResultErr error
		)
		if len(args) == compareNArgs {
			version := args[0]
			otherVersion := args[1]
			compareResult, compareResultErr = gosemver.CompareSemVer(version, otherVersion)
		}
		if len(args) == compareNArgsStdin {
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
			if len(versionsSlice) != compareNArgs {
				fmt.Fprintln(os.Stderr, "Error: two versions should be provided")
				os.Exit(c.ExitOtherErrors)
			}

			compareResult, compareResultErr = gosemver.CompareSemVer(versionsSlice[0], versionsSlice[1])
		}
		if compareResultErr != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", compareResultErr)
			if errors.Is(compareResultErr, gosemver.ErrInvalidVersion) {
				os.Exit(c.ExitInvalidSemver)
			}
			os.Exit(c.ExitOtherErrors)
		}
		fmt.Println(compareResult)
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
