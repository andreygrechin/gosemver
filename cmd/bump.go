package cmd

import (
	"errors"
	"fmt"
	"os"

	c "github.com/andreygrechin/gosemver/internal/config"
	"github.com/andreygrechin/gosemver/pkg/gosemver"
	"github.com/spf13/cobra"
)

var (
	newPrereleaseID string
	newBuildID      string
)

var bumpCmd = &cobra.Command{
	Use:   "bump <semver_id> <version|->",
	Short: "Increment a specific SemVer identifier",
	Long: `Increment specific SemVer identifier <semver_id> of a provided semantic version <version> where
identifier is (major|minor|patch|prerelease|build|release).

The version can be provided either as an argument or via stdin when using '-' as the argument.
Only one input method can be used at a time.

Examples:
  gosemver bump major 0.1.2
  gosemver bump prerelease 2.0.0 --prerelease beta
  gosemver bump prerelease 2.0.0-beta
`,
	Args: cobra.ExactArgs(2), //nolint:mnd
	Run: func(cmd *cobra.Command, args []string) {
		semverID := args[0]
		version, err := gosemver.GetLastArg(*cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting arguments: %v\n", err)
			os.Exit(c.ExitOtherErrors)
		}
		if version == "" {
			fmt.Fprintln(os.Stderr, "Error: empty version string")
			os.Exit(c.ExitOtherErrors)
		}
		if semverID != "prerelease" && newPrereleaseID != "" {
			fmt.Fprintf(os.Stderr, "Error: 'prerelease' flag is allowed only for the 'prerelease' SemVer identifier\n")
			os.Exit(c.ExitOtherErrors)
		}
		if semverID != "build" && newBuildID != "" {
			fmt.Fprintf(os.Stderr, "Error: 'build' flag is allowed only for the 'build' SemVer identifier\n")
			os.Exit(c.ExitOtherErrors)
		}
		semVer, err := gosemver.BumpSemVer(semverID, version, newPrereleaseID, newBuildID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if errors.Is(err, gosemver.ErrInvalidVersion) {
				os.Exit(c.ExitInvalidSemver)
			}
			os.Exit(c.ExitOtherErrors)
		}
		if !gosemver.IsSemVer(semVer.String()) {
			fmt.Fprintf(os.Stderr, "Error: we get an invalid semantic version after bump: %s\n", semVer)
			os.Exit(c.ExitInvalidSemver)
		}
		fmt.Println(semVer)
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.PersistentFlags().StringVarP(
		&newPrereleaseID,
		"prerelease",
		"p",
		"",
		`Add or replace a new prerelease ID, valid only with the 'prerelease' SemVer identifier`,
	)
	bumpCmd.PersistentFlags().StringVarP(
		&newBuildID,
		"build",
		"m",
		"",
		`Add or replace a new build metadata ID, valid only with the 'build' SemVer identifier`,
	)
}
