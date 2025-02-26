package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
	"gopkg.in/yaml.v3"
)

const configFileName = "igittconfig.yaml"

type IgittConfig struct {
	IconType        string `yaml:"iconType"`
	ShowAllCommands bool   `yaml:"showAllCommands"`
}

func InitialConfig() (bool, error) {
	configExists, configPath := HasConfigFile()

	if configExists {
		return false, nil
	}

	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		logger.ErrorLogger.Panic(err)
		return false, err
	}

	_, err = file.Write([]byte(GetDefaultConfig()))

	if err != nil {
		logger.ErrorLogger.Panic(err)
		return false, err
	}

	defer file.Close()

	return true, nil
}

func HasConfigFile() (bool, string) {
	executable, err := os.Executable()
	executableDir := filepath.Dir(executable)
	configPath := filepath.Join(executableDir, configFileName)

	if err != nil {
		logger.ErrorLogger.Panic(err)
	}

	_, err = os.Stat(configPath)
	fileExists := !os.IsNotExist(err)
	isFileEmpty := true

	if fileExists {
		file, err := os.ReadFile(configPath)
		if err != nil {
			logger.ErrorLogger.Panic(err)
		}
		fileContents := string(file)
		isFileEmpty = len(fileContents) == 0
	}

	result := fileExists && !isFileEmpty

	return result, configPath
}

func GetConfig() IgittConfig {
	configExists, configPath := HasConfigFile()
	if !configExists {
		logger.ErrorLogger.Print("Config file does not exist, creating one now")

		_, err := InitialConfig()

		if err != nil {
			logger.ErrorLogger.Panic(err)
			return IgittConfig{}
		}
	}

	config, err := ReadConfigFromPath(configPath)
	if err != nil {
		utilities.PrintGeneralError(fmt.Sprintf("Failed to read the configuration at:\n%s.\n\nYou can either fix it or delete it to generate a new configuration file.", configPath))
		logger.ErrorLogger.Fatal("Failed to read the configuration. Perhaps it is invalid?")
	}

	return config
}

func ReadConfigFromPath(configPath string) (IgittConfig, error) {
	var config IgittConfig
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func GetDefaultConfig() string {
	configContent := `# This is the configuration for Igitt.
# Please adjust the values as needed.

# Choices: "emoji", "unicode", "nerdfont", "ascii" - Default: "unicode"
iconType: "unicode"

# Show all commands in the interactive mode, even if not in a Git repository.
# Default: false
showAllCommands: false
`

	return configContent
}

func GetConfigPath(print bool) string {
	configExists, configPath := HasConfigFile()

	if !configExists {
		logger.ErrorLogger.Print("Config file does not exist, creating one now")
		_, err := InitialConfig()
		if err != nil {
			logger.ErrorLogger.Panic(err)
			return ""
		}
	}

	if print {
		fmt.Printf("\n\nTo edit the configuration, %s in your text editor:\n\n", color.YellowString("open the following file"))
		color.Blue(configPath)
		fmt.Println()
	}

	return configPath
}
