package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	input = strings.TrimRight(input, "\r\n")
	return strings.TrimSpace(input)
}

func ReadLineWithDefault(prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	input := ReadLine()
	if input == "" {
		return defaultValue
	}
	return input
}

func ReadBool(prompt string, defaultValue bool) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	input := strings.ToLower(ReadLine())
	if input == "" {
		return defaultValue
	}
	return input == "y" || input == "yes"
}
