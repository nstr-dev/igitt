package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func AddChanges(arguments string) {
	fmt.Println("Adding changes")

	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "add", arguments).CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error adding changes:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Adding changes:", errOut, string(byteOut))
}
