package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	input = strings.TrimRight(input, "\r\n")
	return strings.TrimSpace(input)
}

func readLineWithDefault(prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	input := readLine()
	if input == "" {
		return defaultValue
	}
	return input
}
