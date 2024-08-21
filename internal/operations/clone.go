package operations

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CloneRepository(repoUrl string) {
	fmt.Println("Cloning repository from " + repoUrl)
	out, out2 := exec.Command("git", "clone", repoUrl).Output()
	logger.InfoLogger.Println(out, out2)
}
