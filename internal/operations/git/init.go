package git

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
)

func InitRepository() {
	mydir, err := os.Getwd()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	fmt.Println("Initializing repository in " + mydir)

	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "init").CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error initializing Git repository:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Initializing Git repository:", errOut, string(byteOut))
}
