package interactive

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahstreller/igitt/internal/utilities/logger"
	"github.com/rivo/uniseg"
	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed operations.json
var commandJSON []byte

const iconWidth = 3

type Command struct {
	Id            string `json:"id"`
	Icon          string `json:"icon"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	NextStep      string `json:"nextStep"`
	NextStepTitle string `json:"nextStepTitle"`
}

func getTitle(command Command) string {
	return command.Icon + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.Icon)) + command.Name
}

func StartInteractive(rootCmd *cobra.Command) {
	var commands []Command
	var selectedCommand Command

	var repoUrlInput string

	_ = json.Unmarshal(commandJSON, &commands)
	commandOptions := make([]huh.Option[Command], len(commands))

	for i, c := range commands {
		title := " " + getTitle(c)
		commandOptions[i] = huh.NewOption(title, c)
	}

	theme := huh.ThemeCatppuccin()
	theme.Focused.Base.Border(lipgloss.HiddenBorder())

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Command]().
				Title("Igitt").
				Height(20).
				DescriptionFunc(func() string {
					if selectedCommand.NextStep != "none" {
						return "\n" + getTitle(selectedCommand) + ": " + selectedCommand.Description + "\n\n" + "     " + "Next step: " + selectedCommand.NextStepTitle + "\n"
					}
					return "\n" + getTitle(selectedCommand) + ": " + selectedCommand.Description + "\n\n\n"
				}, &selectedCommand).
				Options(commandOptions...).
				Value(&selectedCommand),
		),
		huh.NewGroup(huh.NewInput().
			Title("Link to Git repository").
			Description("\n  Enter the link to your repository here.\n").
			Suggestions([]string{"https://github.com"}).
			Value(&repoUrlInput),
		).WithHideFunc(func() bool {
			return selectedCommand.NextStep != "op-enter-repo-url"
		}),
	).WithTheme(theme).Run()

	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	// rootCmd.Print(selectedCommand)
}
