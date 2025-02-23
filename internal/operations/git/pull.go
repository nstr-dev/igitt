package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
)

func PullRemote() {
	fmt.Println("Pulling from remote repository")
	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "pull").CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error pulling from remote repository:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Pulling from remote repository:", errOut, string(byteOut))
}
