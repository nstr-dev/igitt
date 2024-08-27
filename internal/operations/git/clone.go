package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CloneRepository(repoUrl string) {
	fmt.Println("Cloning repository from " + repoUrl)
	byteOut, errOut := exec.Command("git", "clone", repoUrl).Output()
	logger.InfoLogger.Println("Cloning:", errOut, byteOut)
}
