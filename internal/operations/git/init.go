package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func InitRepository() {
	mydir, err := os.Getwd()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	fmt.Println("Initializing repository in" + mydir)
	byteOut, errOut := exec.Command("git", "init").CombinedOutput()
	logger.InfoLogger.Println("Initializing Git repository:", errOut, string(byteOut))
}
