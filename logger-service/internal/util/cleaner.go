package util

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// The cleanString function simply replaces all occurrences of newline and carriage return characters with an empty string.
func cleanString(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	return input
}

// ReadAndCleanString function reads a string from the standard input (os.Stdin) and removes any newline characters (\n) or carriage
// return characters (\r) from the string before returning it.
func ReadAndCleanString() string {
	result, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	result = cleanString(result)
	return result
}
