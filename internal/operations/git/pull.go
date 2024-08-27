package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func PullRemote() {
	fmt.Println("Pulling from remote repository")
	byteOut, errOut := exec.Command("git", "pull").Output()
	logger.InfoLogger.Println("Pulling from remote repository:", errOut, byteOut)
}
