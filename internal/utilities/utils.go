package utilities

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
)

var spacing string = "=============================================================================\n\n"

func RemoveLastEmptyLine(text string) string {
	lines := strings.Split(text, "\n")

	// Remove the last empty line if it exists
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

func PrintGeneralError(message string) {
	fmt.Printf("\n%s", color.HiBlackString(spacing))
	fmt.Printf("%s  There was an issue:\n\n%s\n", color.HiRedString("⚠"), color.HiRedString(message))
	fmt.Printf("\n%s", color.HiBlackString(spacing))
}

func PrintGitError(message string) {
	fmt.Printf("\n%s", color.HiBlackString(spacing))
	fmt.Printf("%s  There was an issue. Received following message from Git:\n\n%s\n", color.HiRedString("⚠"), color.HiRedString(message))
	fmt.Printf("\n%s", color.HiBlackString(spacing))
}

func CheckIsRepo() bool {
	byteOut, errOut := exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()

	if errOut != nil {
		if strings.Contains(string(byteOut), "not a git repository") {
			return false
		}
		logger.ErrorLogger.Println("Error checking if inside a Git repository:", errOut, string(byteOut))
		PrintGitError(string(byteOut))
		return false
	}

	if strings.Contains(string(byteOut), "true") {
		return true
	}

	return false
}
