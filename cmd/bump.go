package cmd

import (
	"fmt"
	"os"

	"github.com/andreygrechin/gosemver/pkg/semver"
	"github.com/spf13/cobra"
)

var newPrereleaseID string

var bumpCmd = &cobra.Command{
	Use:   "bump <semver_id> <version>",
	Short: "Increment a specific SemVer identifier",
	Long: `Increment specific SemVer identifier <semver_id> of a provided semantic version <version> where
identifier is (major|minor|patch|prerelease|release).

'prerelease' identifier may be specified with a prerelease ID flag.

Examples:
  gosemver bump major 0.1.2
  gosemver bump prerelease 2.0.0 --prerelease-id beta
  gosemver bump prerelease 2.0.0-beta
`,
	Args: cobra.ExactArgs(2), //nolint:mnd
	Run: func(cmd *cobra.Command, args []string) {
		semverID := args[0]
		version := args[1]
		if semverID != "prerelease" && newPrereleaseID != "" {
			fmt.Printf("error: 'prerelease-id' flag is allowed only for the 'prerelease' SemVer identifier\n")
			os.Exit(1)
		}
		semVer, err := semver.BumpSemVer(semverID, version, newPrereleaseID)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		if !semver.IsSemVer(semVer.String()) {
			fmt.Printf("error: we get an invalid semantic version after bump: %s\n", semVer)
			os.Exit(1)
		}
		fmt.Println(semVer)
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.PersistentFlags().StringVarP(
		&newPrereleaseID,
		"prerelease-id",
		"p",
		"",
		`A new prerelease ID to add or replace the existing one, valid only for the
'prerelease' SemVer identifier`,
	)
}
