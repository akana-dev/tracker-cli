package updater

import (
	"tracker/internal/version"
)

func isDevVersion() bool {
	return version.IsDev()
}

func getCurrentVersion() string {
	return version.Version
}
