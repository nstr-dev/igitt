package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CommitChanges(message string) {
	fmt.Println("Committing changes")
	byteOut, errOut := exec.Command("git", "commit", "-m", message).CombinedOutput()

	if errOut != nil {
		logger.ErrorLogger.Println("Error committing changes:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Committing changes:", errOut, string(byteOut))
}
