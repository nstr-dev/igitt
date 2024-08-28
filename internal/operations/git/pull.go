package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func PullRemote() {
	fmt.Println("Pulling from remote repository")
	byteOut, errOut := exec.Command("git", "pull").CombinedOutput()

	if errOut != nil {
		logger.ErrorLogger.Println("Error pulling from remote repository:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Pulling from remote repository:", errOut, string(byteOut))
}
