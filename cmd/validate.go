package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [version|-]",
	Short: "Validate a semantic version",
	Long: `Validate if a provided semantic version string complies with semver 2.0.0 specification. Exits with
status 0 if valid, 1 if invalid. Prints "valid" or "invalid" to stdout.

The version can be provided either as an argument or via stdin when using '-' as the argument.
Only one input method can be used at a time.

Examples:
  gosemver validate 1.2.3
  gosemver validate v1.2.3-beta.1+build.123
  echo "1.2.3" | gosemver validate -
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var version string

		if args[0] == "-" {
			reader := bufio.NewReader(cmd.InOrStdin())
			input, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
				os.Exit(1)
			}
			version = strings.TrimSpace(input)
		} else {
			version = args[0]
		}

		if version == "" {
			fmt.Fprintln(os.Stderr, "Error: empty version string")
			os.Exit(1)
		}

		if gosemver.IsSemVer(version) {
			fmt.Println("valid")
		} else {
			fmt.Println("invalid")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
