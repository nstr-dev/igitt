package operations

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func CreateAliasScripts() {
	executable, err := os.Executable()
	if err != nil {
		logger.ErrorLogger.Println("Failed to get executable path:", err)
		return
	}
	executablePath := filepath.Dir(executable)
	logger.InfoLogger.Printf("Creating alias for executable in %v", executablePath)

	aliasScriptWindows := `@ECHO OFF
if exist `+ executable +` (
	`+ executable +` %*
) else (
	igitt mkalias
	igitt %*
)`

	aliasScriptUnix := `#!/bin/bash
"` + executable + `" "$@"`

	var aliasFileName, aliasScriptContent string

	switch runtime.GOOS {
	case "windows":
		aliasFileName = filepath.Join(executablePath, "igt.cmd")
		aliasScriptContent = aliasScriptWindows
	case "linux", "darwin":
		aliasFileName = filepath.Join(executablePath, "igt")
		aliasScriptContent = aliasScriptUnix
	default:
		logger.ErrorLogger.Println("Unsupported OS")
		return
	}

	if err := os.WriteFile(aliasFileName, []byte(aliasScriptContent), 0755); err != nil {
		logger.ErrorLogger.Println("Failed to create alias script:", err)
		return
	}

	fmt.Printf("Alias script created successfully at %v\n", color.GreenString(aliasFileName))
	logger.InfoLogger.Printf("Alias script created successfully at %v", aliasFileName)
}
