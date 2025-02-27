package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
)

func PushRemote() {
	fmt.Println("Pushing to remote repository")
	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "-c", "push.autoSetupRemote=true", "push").CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error pushing to remote repository:", errOut, string(byteOut))
		utilities.PrintGitError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Pushing to remote repository:", errOut, string(byteOut))
}
