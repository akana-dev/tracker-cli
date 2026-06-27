package main

import (
	"errors"
	"os"

	"github.com/fatih/color"

	"tracker/internal/aliases"
	"tracker/internal/cli"
)

func main() {
	os.Args = aliases.ExpandArgs(os.Args)

	if err := cli.Execute(); err != nil {
		if errors.Is(err, cli.ErrHelp) {
			os.Exit(0)
		}

		errPrinter := color.New(color.FgRed, color.Bold)
		errPrinter.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
