package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/noahstreller/igitt/internal/utilities/logger"
	"gopkg.in/yaml.v3"
)

const configFileName = "igittconfig.yaml"

type IgittConfig struct {
	IconType string `yaml:"iconType"`
}

func InitialConfig() error {

	configExists, configPath := HasConfigFile()
	if configExists {
		return nil
	}

	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return err
	}

	_, err = file.Write([]byte(GetDefaultConfig()))

	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return err
	}

	defer file.Close()

	return nil
}

func HasConfigFile() (bool, string) {
	executable, err := os.Executable()
	executableDir := filepath.Dir(executable)
	configPath := filepath.Join(executableDir, configFileName)

	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	_, err = os.Stat(configPath)

	return !os.IsNotExist(err), configPath
}

func GetConfig() IgittConfig {
	configExists, configPath := HasConfigFile()
	if !configExists {
		logger.ErrorLogger.Fatal("Config file does not exist, creating one now")
		err := InitialConfig()
		if err != nil {
			logger.ErrorLogger.Fatal(err)
			return IgittConfig{}
		}
	}

	config, err := ReadConfigFromPath(configPath)
	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return IgittConfig{}
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

# Choices: "emoji", "unicode", "nerdfont"
iconType: "unicode"
`

	return configContent
}

func GetConfigPath(print bool) string {
	configExists, configPath := HasConfigFile()

	if !configExists {
		logger.ErrorLogger.Fatal("Config file does not exist, creating one now")
		err := InitialConfig()
		if err != nil {
			logger.ErrorLogger.Fatal(err)
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
