package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func AddChanges(arguments string) {
	fmt.Println("Adding changes")

	byteOut, errOut := exec.Command("git", "add", arguments).CombinedOutput()
	if errOut != nil {
		logger.ErrorLogger.Println("Error adding changes:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}
	logger.InfoLogger.Println("Adding changes:", errOut, string(byteOut))
}
