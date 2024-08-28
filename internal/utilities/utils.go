package utilities

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func RemoveLastEmptyLine(text string) string {
	lines := strings.Split(text, "\n")

	// Remove the last empty line if it exists
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

func PrintError(message string) {
	fmt.Printf("\n%s", color.HiBlackString("=============================================================================\n\n"))
	fmt.Printf("%s  There was an issue. Received following message from Git:\n\n%s\n", color.HiRedString("⚠"), color.HiRedString(message))
}
