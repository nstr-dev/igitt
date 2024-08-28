package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CloneRepository(repoUrl string) {
	fmt.Println("Cloning repository from " + repoUrl)

	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "clone", repoUrl).CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error cloning:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}

	logger.InfoLogger.Println("Cloning:", errOut, string(byteOut))
}
