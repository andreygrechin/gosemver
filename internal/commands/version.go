package commands

import (
	"fmt"

	"github.com/andreygrechin/gosemver/internal/config"
)

func Version(args []string) {
	fmt.Printf("Version: %s\nBuild Time: %s\nCommit: %s\n", config.Version, config.BuildTime, config.Commit)
}
