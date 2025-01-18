package gosemver

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// bumpExistingNumeric tries to bump the numeric suffix in an existing prerelease.
func bumpExistingNumeric(existing string) string {
	prefix, numeric := splitNumericSuffix(existing)
	if numeric != "" {
		oldNum, _ := strconv.Atoi(numeric)
		oldNum++
		return fmt.Sprintf("%s%d", prefix, oldNum)
	}
	// if there's no numeric part, append "1"
	return fmt.Sprintf("%s1", prefix)
}

// splitNumericSuffix extracts the prefix part (could include letters, dots, hyphens)
// and the trailing numeric part if any.
func splitNumericSuffix(prerelease string) (string, string) {
	idx := -1

	for i := len(prerelease) - 1; i >= 0; i-- {
		if prerelease[i] < '0' || prerelease[i] > '9' {
			// first non-digit from the end => i+1 is the start of numeric
			idx = i

			break
		}
	}

	if idx == len(prerelease)-1 {
		// no trailing digits
		return prerelease, ""
	}
	return prerelease[:idx+1], prerelease[idx+1:]
}

func GetLastArg(cmd cobra.Command, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no arguments provided")
	}
	if args[len(args)-1] == "-" {
		reader := bufio.NewReader(cmd.InOrStdin())
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("error reading from stdin: %w", err)
		}
		return strings.TrimSpace(input), nil
	}
	return args[len(args)-1], nil
}
