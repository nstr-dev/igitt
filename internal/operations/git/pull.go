package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func PullRemote() {
	fmt.Println("Pulling from remote repository")
	byteOut, errOut := exec.Command("git", "pull").CombinedOutput()
	logger.InfoLogger.Println("Pulling from remote repository:", errOut, string(byteOut))
}
