package cli

import (
	"tracker/internal/input"
)

func readLine() string {
	return input.ReadLine()
}

func readLineWithDefault(prompt, defaultValue string) string {
	return input.ReadLineWithDefault(prompt, defaultValue)
}
