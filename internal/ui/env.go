package ui

import (
	"os"
	"sync"

	"github.com/mattn/go-isatty"
)

var (
	PlainMode bool
	noColor   bool
	isTTY     bool
	once      sync.Once
)

func init() {
	once.Do(func() {
		isTTY = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
		noColor = os.Getenv("NO_COLOR") != ""

		if os.Getenv("TRACKER_PLAIN") != "" {
			PlainMode = true
		}
	})
}

func IsInteractive() bool {
	return isTTY && !noColor && !PlainMode
}

func ShouldUseColors() bool {
	return isTTY && !noColor && !PlainMode
}

func IsTTY() bool {
	return isTTY
}
