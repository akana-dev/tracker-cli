package version

import "fmt"

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("tracker version %s (commit: %s, built: %s)", Version, Commit, BuildDate)
}

func IsDev() bool {
	return Version == "dev" || Version == ""
}
