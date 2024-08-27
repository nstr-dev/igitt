package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CommitChanges(message string) {
	fmt.Println("Committing changes")
	byteOut, errOut := exec.Command("git", "commit", "-m", message).CombinedOutput()
	logger.InfoLogger.Println("Committing changes:", errOut, string(byteOut))
}
