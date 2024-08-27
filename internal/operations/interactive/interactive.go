package interactive

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/noahstreller/igitt/internal/operations/git"
	"github.com/noahstreller/igitt/internal/utilities/logger"
	"github.com/rivo/uniseg"

	_ "embed"
)

//go:embed operations.json
var commandJSON []byte

const (
	Unicode IconType = iota
	Emoji
	NerdFont
)

type IconType int

func (d IconType) String() string {
	return [...]string{"Unicode", "Emoji", "NerdFont"}[d]
}

type Command struct {
	Id            string `json:"id"`
	Icon          string `json:"icon"`
	IconEmoji     string `json:"icon_emoji"`
	IconNerdFont  string `json:"icon_nerdfont"`
	Name          string `json:"name"`
	Shortcut      string `json:"shortcut"`
	Description   string `json:"description"`
	NextStep      string `json:"nextStep"`
	NextStepTitle string `json:"nextStepTitle"`
}

type CommandFlowResult struct {
	SelectedCommand Command
	RepoUrlInput    string
	GitAddArguments string
}

const iconWidth = 3
const shortcutsEnabled = false
const iconVariant = NerdFont

func getTitle(command Command) string {
	if iconVariant == Emoji {
		return command.IconEmoji + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.IconEmoji)) + command.Name
	}
	if iconVariant == NerdFont {
		return command.IconNerdFont + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.IconNerdFont)) + command.Name
	}
	return command.Icon + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.Icon)) + command.Name
}

func getNextStepIcon(variant IconType) string {
	if variant == Emoji {
		return "⏩"
	}
	if variant == NerdFont {
		return ""
	}
	return "↪"
}

func getNoNextStepIcon(variant IconType) string {
	if variant == Emoji {
		return "🎯"
	}
	if variant == NerdFont {
		return ""
	}
	return "◎"
}

func getLinkIcon(variant IconType) string {
	if variant == Emoji {
		return "🔗  "
	}
	if variant == NerdFont {
		return "  "
	}
	return ""
}

func StartInteractive() {
	var commands []Command

	commandFlowResult := CommandFlowResult{
		SelectedCommand: Command{Id: "none"},
		RepoUrlInput:    "",
		GitAddArguments: "",
	}

	_ = json.Unmarshal(commandJSON, &commands)
	commandOptions := make([]huh.Option[Command], len(commands))

	for i, c := range commands {
		if c.Shortcut == "none" || !shortcutsEnabled {
			formattedTitle := " " + getTitle(c)
			commandOptions[i] = huh.NewOption(formattedTitle, c)
		} else {
			formattedTitle := fmt.Sprintf("%-15s(%s)", " "+getTitle(c), c.Shortcut)
			commandOptions[i] = huh.NewOption(formattedTitle, c)
		}
	}

	theme := huh.ThemeCatppuccin()
	theme.Focused.Base.Border(lipgloss.HiddenBorder())
	theme.Form.Border(lipgloss.NormalBorder())

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Command]().
				Title("Igitt").
				Filtering(true).
				Height(20).
				DescriptionFunc(func() string {
					if commandFlowResult.SelectedCommand.NextStep != "none" {
						return "\n" + getTitle(commandFlowResult.SelectedCommand) + ": " + commandFlowResult.SelectedCommand.Description + "\n\n" + "   " + getNextStepIcon(iconVariant) + "  " + "Next step: " + commandFlowResult.SelectedCommand.NextStepTitle + "\n"
					}
					return "\n" + getTitle(commandFlowResult.SelectedCommand) + ": " + commandFlowResult.SelectedCommand.Description + "\n\n" + "   " + getNoNextStepIcon(iconVariant) + "  " + "No next steps" + "\n"
				}, &commandFlowResult.SelectedCommand).
				Options(commandOptions...).
				Value(&commandFlowResult.SelectedCommand),
		),

		huh.NewGroup(huh.NewInput().
			Title("Link to Git repository").
			Description("\n"+getLinkIcon(iconVariant)+"Enter the link to your repository here.\n").
			Suggestions([]string{
				"https://github.com/",
				"https://gitlab.com/",
				"https://bitbucket.org/",
				"https://sourcehut.org/",
				"https://codeberg.org/",
				"https://gitea.io/",
			}).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("if you're ready to clone, enter a repository URL")
				}
				return nil
			}).
			Value(&commandFlowResult.RepoUrlInput),
		).WithHideFunc(func() bool {
			return commandFlowResult.SelectedCommand.NextStep != "op-enter-repo-url"
		}),
	).WithTheme(theme).Run()

	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	runResultingCommand(commandFlowResult)
}

func runResultingCommand(commandFlow CommandFlowResult) {
	if commandFlow.SelectedCommand.Id == "op-clone" && commandFlow.RepoUrlInput != "" {
		logger.InfoLogger.Println("clone command selected, sending to operations")
		git.CloneRepository(commandFlow.RepoUrlInput)
		return
	}

	if commandFlow.SelectedCommand.Id == "op-init" {
		logger.InfoLogger.Println("init command selected, sending to operations")
		git.InitRepository()
		return
	}

	if commandFlow.SelectedCommand.Id == "op-pull" {
		logger.InfoLogger.Println("pull command selected, sending to operations")
		git.PullRemote()
		return
	}

	if commandFlow.SelectedCommand.Id == "op-push" {
		logger.InfoLogger.Println("push command selected, sending to operations")
		git.PushRemote()
		return
	}

	if commandFlow.SelectedCommand.Id == "op-add" && commandFlow.GitAddArguments != "" {
		logger.InfoLogger.Println("add command selected, sending to operations")
		git.AddChanges(commandFlow.GitAddArguments)
		return
	}
}
