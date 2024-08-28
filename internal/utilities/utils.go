package utilities

import "strings"

func RemoveLastEmptyLine(text string) string {
	lines := strings.Split(text, "\n")

	// Remove the last empty line if it exists
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}
