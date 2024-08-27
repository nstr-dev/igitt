package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func PushRemote() {
	fmt.Println("Pushing to remote repository")
	byteOut, errOut := exec.Command("git", "push").Output()
	logger.InfoLogger.Println("Pushing to remote repository:", errOut, byteOut)
}
