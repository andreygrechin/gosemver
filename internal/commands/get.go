package commands

import (
	"errors"
	"fmt"

	"github.com/andreygrechin/gosemver/internal/semver"
)

func GetSemVer(subCommand, version string) (*semver.SemVer, error) {
	ver, err := semver.ParseSemVer(version)
	if err != nil {
		return nil, err
	}
	switch subCommand {
	case "major":
		fmt.Println(ver.Major)
		return &semver.SemVer{Major: ver.Major}, nil
	case "minor":
		fmt.Println(ver.Minor)
		return &semver.SemVer{Minor: ver.Minor}, nil
	case "patch":
		fmt.Println(ver.Patch)
		return &semver.SemVer{Patch: ver.Patch}, nil
	case "prerel":
		fmt.Println(ver.Prerelease)
		return &semver.SemVer{Prerelease: ver.Prerelease}, nil
	case "prerelease":
		fmt.Println(ver.Prerelease)
		return &semver.SemVer{Prerelease: ver.Prerelease}, nil
	case "release":
		fmt.Println(ver.Release)
		return &semver.SemVer{Release: ver.Release}, nil
	case "build":
		fmt.Println(ver.Build)
		return &semver.SemVer{Build: ver.Build}, nil
	default:
		return nil, errors.New("unknown get command")
	}
}
