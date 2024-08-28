package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CloneRepository(repoUrl string) {
	fmt.Println("Cloning repository from " + repoUrl)
	byteOut, errOut := exec.Command("git", "clone", repoUrl).CombinedOutput()

	if errOut != nil {
		logger.ErrorLogger.Println("Error cloning:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}

	logger.InfoLogger.Println("Cloning:", errOut, string(byteOut))
}
