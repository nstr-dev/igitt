package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
)

func CommitChanges(message string) {
	fmt.Println("Committing changes")
	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "commit", "-m", message).CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		utilities.PrintGitError(string(byteOut))
		logger.ErrorLogger.Fatal("Error committing changes:", errOut, string(byteOut))
		return
	}
	logger.InfoLogger.Println("Committing changes:", errOut, string(byteOut))
}
