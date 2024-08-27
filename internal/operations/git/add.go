package git

import (
	"fmt"
	"os/exec"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func AddChanges(arguments string) {
	fmt.Println("Adding changes")

	byteOut, errOut := exec.Command("git", "add", arguments).CombinedOutput()
	logger.InfoLogger.Println("Adding changes:", errOut, string(byteOut))
}
