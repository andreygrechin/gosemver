package config

const (
	ExitOK            = 0
	ExitInvalidSemver = 1
	ExitOtherErrors   = 2
)

var (
	Version   = "unknown"
	BuildTime = "unknown"
	Commit    = "unknown"
)
