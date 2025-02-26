package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
)

func AddEverything() {
	AddChanges(".")
}

func AddChanges(arguments string) {
	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "add", arguments).CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error adding changes:", errOut, string(byteOut))
		utilities.PrintGitError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Adding changes:", errOut, string(byteOut))

	fmt.Println("Changes added to the staging area.")
}
