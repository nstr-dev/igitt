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
	NewBranchName   string
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
			b = fmt.Sprintf("%s *", b)
		}
		branchOptions[i] = huh.NewOption(b, b)
	}

	branchOptions = append(branchOptions, huh.NewOption("[ Create new branch ]", "[newBranch]"))

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

var commandFlowResult = CommandFlowResult{
	SelectedCommand: Command{Id: "none"},
	RepoUrlInput:    "",
	GitAddArguments: "",
	CommitMessage:   "",
	NewBranchName:   "",
	SelectedBranch:  "",
	BranchAction:    "",
}

func StartInteractive() {
	var commands []Command
	formGroups := make(map[string]*huh.Form)

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

	formGroups["ns-choose-branch"] =
		huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Branch selection").
					Description(bold("\n  Select a branch\n")).
					OptionsFunc(getBranchOptions, &commandFlowResult).
					Value(&commandFlowResult.SelectedBranch))).WithTheme(theme)

	formGroups["ns-choose-branch-action"] =
		huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Branch action selection").
					Description(fmt.Sprintf("\n  Select an action for the branch %s\n", commandFlowResult.SelectedBranch)).
					Options(getBranchActionOptions()...).
					Value(&commandFlowResult.BranchAction))).WithTheme(theme)

	formGroups["ns-enter-new-branch-name"] =
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Branch name").
					Description("\n" + "Enter the desired branch name here.\n").
					Suggestions([]string{
						"feature/",
						"bugfix/",
						"hotfix/",
						"fix/",
						"refactor/",
						"chore/",
						"docs/",
					}).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("the branch name should not be empty")
						}

						if strings.ContainsAny(s, " ~^:?*[]\\") ||
							strings.Contains(s, "\\") ||
							strings.Contains(s, "//") ||
							strings.Contains(s, "@{") ||
							strings.Contains(s, "..") ||
							s == "@" {
							return fmt.Errorf("some special characters are not allowed in branch names")
						}

						return nil
					}).
					Value(&commandFlowResult.NewBranchName))).WithTheme(theme)

	formGroups["ns-enter-repo-url"] =
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Link to Git repository").
					Description("\n" + getLinkIcon(iconVariant) + "Enter the link to your repository here.\n").
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
					Value(&commandFlowResult.RepoUrlInput))).WithTheme(theme)

	formGroups["ns-enter-commit-message"] =
		huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Commit message").
					Description("\n" + getCommitIcon(iconVariant) + "Type a short description to the commit.\n").
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
					Value(&commandFlowResult.CommitMessage))).WithTheme(theme)

	mainForm := huh.NewForm(
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
	).WithTheme(theme).WithHeight(len(commands) + 9)

	err := mainForm.Run()
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	nextStepErr := runNextStep(formGroups)
	if nextStepErr != nil {
		logger.ErrorLogger.Fatal(nextStepErr)
	} else {
		runResultingCommand()
	}
}

func runNextStep(formGroups map[string]*huh.Form) error {
	if commandFlowResult.SelectedCommand.NextStep == "ns-choose-branch" {
		err := formGroups["ns-choose-branch"].Run()
		if err != nil {
			return err
		}

		fmt.Println("Selected branch:", commandFlowResult.SelectedBranch)

		if commandFlowResult.SelectedBranch == "[newBranch]" {
			fmt.Println("yes")
			commandFlowResult.SelectedCommand.NextStep = "ns-enter-new-branch-name"
			err = runNextStep(formGroups)
			return err
		}

		fmt.Println("no")
		commandFlowResult.SelectedCommand.NextStep = "ns-choose-branch-action"
		err = runNextStep(formGroups)
		return err
	}

	if commandFlowResult.SelectedCommand.NextStep == "ns-enter-new-branch-name" {
		return formGroups["ns-enter-new-branch-name"].Run()
	}

	if commandFlowResult.SelectedCommand.NextStep == "ns-choose-branch-action" {
		return formGroups["ns-choose-branch-action"].Run()
	}

	if commandFlowResult.SelectedCommand.NextStep == "ns-enter-repo-url" {
		return formGroups["ns-enter-repo-url"].Run()
	}

	if commandFlowResult.SelectedCommand.NextStep == "ns-enter-commit-message" {
		return formGroups["ns-enter-commit-message"].Run()
	}

	return nil
}

func runResultingCommand() {

	if commandFlowResult.SelectedCommand.Id == "op-clone" && commandFlowResult.RepoUrlInput != "" {
		logger.InfoLogger.Println("clone command selected, sending to operations")
		git.CloneRepository(commandFlowResult.RepoUrlInput)
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-commit" && commandFlowResult.CommitMessage != "" {
		logger.InfoLogger.Println("commit command selected, sending to operations")
		git.CommitChanges(commandFlowResult.CommitMessage)
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-branches" &&
		commandFlowResult.SelectedBranch != "" &&
		commandFlowResult.NewBranchName != "" &&
		commandFlowResult.SelectedCommand.NextStep == "ns-enter-new-branch-name" {

		logger.InfoLogger.Printf("branch command selected, sending to operations, new branch name: %s\n", commandFlowResult.NewBranchName)
		git.CreateBranch(commandFlowResult.NewBranchName)
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-branches" &&
		commandFlowResult.SelectedBranch != "" &&
		commandFlowResult.BranchAction != "" {

		isCheckedOutAlready := strings.Contains(commandFlowResult.SelectedBranch, "*")
		checkedOutBranchWithoutStar := strings.TrimPrefix(commandFlowResult.SelectedBranch, "*")

		logger.InfoLogger.Printf("branch command selected, sending to action block, isCheckedOutAlready: %v, checkedOutBranchWithoutStar: %s\n", isCheckedOutAlready, checkedOutBranchWithoutStar)

		if commandFlowResult.BranchAction == "Check out" && !isCheckedOutAlready {
			logger.InfoLogger.Println("checkout command selected, sending to operations")
			git.CheckoutBranch(commandFlowResult.SelectedBranch)
			return
		}
		if commandFlowResult.BranchAction == "Check out" && isCheckedOutAlready {
			logger.InfoLogger.Println("checkout command selected, not sending to operations, branch already checked out")
			fmt.Println("Branch already checked out")
			return
		}
	}

	if commandFlowResult.SelectedCommand.Id == "op-init" {
		logger.InfoLogger.Println("init command selected, sending to operations")
		git.InitRepository()
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-status" {
		logger.InfoLogger.Println("status command selected, sending to operations")
		git.Status()
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-pull" {
		logger.InfoLogger.Println("pull command selected, sending to operations")
		git.PullRemote()
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-push" {
		logger.InfoLogger.Println("push command selected, sending to operations")
		git.PushRemote()
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-add" && commandFlowResult.GitAddArguments != "" {
		logger.InfoLogger.Println("add command selected, sending to operations")
		git.AddChanges(commandFlowResult.GitAddArguments)
		return
	}
}
