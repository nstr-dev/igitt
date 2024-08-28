package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func InitRepository() {
	mydir, err := os.Getwd()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	fmt.Println("Initializing repository in" + mydir)
	byteOut, errOut := exec.Command("git", "init").CombinedOutput()

	if errOut != nil {
		logger.ErrorLogger.Println("Error initializing Git repository:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Initializing Git repository:", errOut, string(byteOut))
}
