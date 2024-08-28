package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func PushRemote() {
	fmt.Println("Pushing to remote repository")
	byteOut, errOut := exec.Command("git", "push").CombinedOutput()

	if errOut != nil {
		logger.ErrorLogger.Println("Error pushing to remote repository:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Pushing to remote repository:", errOut, string(byteOut))
}
