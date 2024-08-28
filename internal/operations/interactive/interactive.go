package interactive

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
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
	CommitMessage   string
	SelectedBranch  string
	BranchAction    string
}

const iconWidth = 3
const shortcutsEnabled = false
const iconVariant = NerdFont

func bold(s string) string {
	return color.New(color.Bold).Sprint(s)
}

func getTitle(command Command) string {
	if iconVariant == Emoji {
		return command.IconEmoji + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.IconEmoji)) + bold(command.Name)
	}
	if iconVariant == NerdFont {
		return command.IconNerdFont + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.IconNerdFont)) + bold(command.Name)
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

func getCommitIcon(variant IconType) string {
	if variant == Emoji {
		return "📝  "
	}
	if variant == NerdFont {
		return "  "
	}
	return "✎  "
}

func getBranchOptions() []huh.Option[string] {
	branchResult := git.GetBranches()
	branches := branchResult.Branches

	branchOptions := make([]huh.Option[string], len(branches))

	for i, b := range branches {
		if b == branchResult.CheckedOutBranch {
			b = fmt.Sprintf("%s*", b)
		}
		branchOptions[i] = huh.NewOption(b, b)
	}

	return branchOptions
}

func getBranchActionOptions() []huh.Option[string] {
	branchActions := []string{"Check out", "Delete (wip)", "Rename (wip)"}

	branchActionOptions := make([]huh.Option[string], len(branchActions))

	for i, b := range branchActions {
		branchActionOptions[i] = huh.NewOption(b, b)
	}

	return branchActionOptions
}

func StartInteractive() {
	var commands []Command

	commandFlowResult := CommandFlowResult{
		SelectedCommand: Command{Id: "none"},
		RepoUrlInput:    "",
		GitAddArguments: "",
		CommitMessage:   "",
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
				DescriptionFunc(func() string {
					if commandFlowResult.SelectedCommand.NextStep != "none" {
						return fmt.Sprintf(
							"\n%s: %s\n\n   %s  %s\n\n",
							getTitle(commandFlowResult.SelectedCommand),
							commandFlowResult.SelectedCommand.Description,
							getNextStepIcon(iconVariant),
							bold("Next step: "+commandFlowResult.SelectedCommand.NextStepTitle),
						)
					}

					return fmt.Sprintf(
						"\n%s: %s\n\n   %s  %s next steps\n\n",
						getTitle(commandFlowResult.SelectedCommand),
						commandFlowResult.SelectedCommand.Description,
						getNoNextStepIcon(iconVariant),
						bold("No"),
					)
				}, &commandFlowResult.SelectedCommand).
				Options(commandOptions...).
				Value(&commandFlowResult.SelectedCommand),
		),

		huh.NewGroup(
			huh.NewInput().
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
			return commandFlowResult.SelectedCommand.NextStep != "ns-enter-repo-url"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Commit message").
				Description("\n"+getCommitIcon(iconVariant)+"Type a short description to the commit.\n").
				Suggestions([]string{
					"feat: ",
					"fix: ",
					"docs: ",
					"style: ",
					"refactor: ",
					"perf: ",
					"test: ",
					"chore: ",
				}).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("please enter a commit message")
					}
					return nil
				}).
				Value(&commandFlowResult.CommitMessage),
		).WithHideFunc(func() bool {
			return commandFlowResult.SelectedCommand.NextStep != "ns-enter-commit-message"
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Branch selection").
				Description(bold("\n  Select a branch\n")).
				Options(getBranchOptions()...).
				Value(&commandFlowResult.SelectedBranch),
		).WithHideFunc(func() bool {
			return commandFlowResult.SelectedCommand.NextStep != "ns-choose-branch"
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Branch action selection").
				Description(fmt.Sprintf("\n  Select an action for the branch %s\n", commandFlowResult.SelectedBranch)).
				Options(getBranchActionOptions()...).
				Value(&commandFlowResult.BranchAction),
		).WithHideFunc(func() bool {
			return !(commandFlowResult.SelectedBranch != "")
		}),
	).WithTheme(theme).WithHeight(len(commands) + 9).Run()

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

	if commandFlow.SelectedCommand.Id == "op-commit" && commandFlow.CommitMessage != "" {
		logger.InfoLogger.Println("commit command selected, sending to operations")
		git.CommitChanges(commandFlow.CommitMessage)
		return
	}

	if commandFlow.SelectedCommand.Id == "op-branches" && commandFlow.SelectedBranch != "" && commandFlow.BranchAction != "" {
		isCheckedOutAlready := strings.Contains(commandFlow.SelectedBranch, "*")
		checkedOutBranchWithoutStar := strings.TrimPrefix(commandFlow.SelectedBranch, "*")

		logger.InfoLogger.Printf("branch command selected, sending to action block, isCheckedOutAlready: %v, checkedOutBranchWithoutStar: %s\n", isCheckedOutAlready, checkedOutBranchWithoutStar)

		if commandFlow.BranchAction == "Check out" && !isCheckedOutAlready {
			logger.InfoLogger.Println("checkout command selected, sending to operations")
			git.CheckoutBranch(commandFlow.SelectedBranch)
			return
		}
		if commandFlow.BranchAction == "Check out" && isCheckedOutAlready {
			logger.InfoLogger.Println("checkout command selected, not sending to operations, branch already checked out")
			fmt.Println("Branch already checked out")
			return
		}
	}

	if commandFlow.SelectedCommand.Id == "op-init" {
		logger.InfoLogger.Println("init command selected, sending to operations")
		git.InitRepository()
		return
	}

	if commandFlow.SelectedCommand.Id == "op-status" {
		logger.InfoLogger.Println("status command selected, sending to operations")
		git.Status()
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
