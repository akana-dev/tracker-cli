package main

import (
	"errors"
	"os"

	"github.com/fatih/color"

	"tracker/internal/aliases"
	"tracker/internal/app"
	"tracker/internal/cli"
	"tracker/internal/ui"
)

func main() {
	os.Args = aliases.ExpandArgs(os.Args)

	ui.InitTheme()

	if err := app.InitProvider(); err != nil {
		errPrinter := color.New(color.FgRed, color.Bold)
		errPrinter.Fprintf(os.Stderr, "Ошибка инициализации клиента: %v\n", err)
		os.Exit(1)
	}
	defer app.Cleanup()

	if err := cli.Execute(); err != nil {
		if errors.Is(err, cli.ErrHelp) {
			os.Exit(0)
		}

		errPrinter := color.New(color.FgRed, color.Bold)
		errPrinter.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
