package interactive

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/nstr-dev/igitt/internal/operations/git"
	"github.com/nstr-dev/igitt/internal/utilities"
	"github.com/nstr-dev/igitt/internal/utilities/config"
	"github.com/nstr-dev/igitt/internal/utilities/icons"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
	"github.com/rivo/uniseg"

	_ "embed"
)

//go:embed operations.json
var commandJSON []byte

const (
	Unicode icons.IconType = iota
	Emoji
	NerdFont
	Ascii
)

type Command struct {
	Id              string `json:"id"`
	Icon            string `json:"icon"`
	IconEmoji       string `json:"icon_emoji"`
	IconNerdFont    string `json:"icon_nerdfont"`
	IconAscii       string `json:"icon_ascii"`
	Name            string `json:"name"`
	Shortcut        string `json:"shortcut"`
	Description     string `json:"description"`
	NextStep        string `json:"nextStep"`
	NextStepTitle   string `json:"nextStepTitle"`
	InsideRepoOnly  bool   `json:"insideRepoOnly"`
	OutsideRepoOnly bool   `json:"outsideRepoOnly"`
}

type CommandFlowResult struct {
	SelectedCommand     Command
	RepoUrlInput        string
	GitAddArguments     string
	CommitMessage       string
	SelectedBranch      string
	BranchAction        string
	NewBranchName       string
	DeleteBranchConfirm bool
}

const iconWidth = 3
const shortcutsEnabled = false

func getIconVariantFromConfig() icons.IconType {
	config := config.GetConfig()
	userIconType := strings.ToLower(config.IconType)

	if userIconType == "emoji" {
		return Emoji
	}

	if userIconType == "nerdfont" {
		return NerdFont
	}

	if userIconType == "ascii" {
		return Ascii
	}

	return Unicode
}

func getShowAllCommandsFromConfig() bool {
	config := config.GetConfig()
	return config.ShowAllCommands
}

func getTitle(command Command) string {
	if getIconVariantFromConfig() == Emoji {
		return command.IconEmoji + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.IconEmoji)) + command.Name
	}
	if getIconVariantFromConfig() == NerdFont {
		return command.IconNerdFont + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.IconNerdFont)) + command.Name
	}
	if getIconVariantFromConfig() == Ascii {
		return command.IconAscii + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.Icon)) + command.Name
	}
	return command.Icon + strings.Repeat(" ", iconWidth-uniseg.StringWidth(command.Icon)) + command.Name
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
	branchActions := []string{"Check out", "Delete"}

	branchActionOptions := make([]huh.Option[string], len(branchActions))

	for i, b := range branchActions {
		branchActionOptions[i] = huh.NewOption(b, b)
	}

	return branchActionOptions
}

var commandFlowResult = CommandFlowResult{
	SelectedCommand:     Command{Id: "none"},
	RepoUrlInput:        "",
	GitAddArguments:     "",
	CommitMessage:       "",
	NewBranchName:       "",
	SelectedBranch:      "",
	BranchAction:        "",
	DeleteBranchConfirm: false,
}

func StartInteractive() {

	interactiveTitle := color.New(color.Bold, color.FgGreen).PrintfFunc()
	var interactiveTitleText string
	var interactiveByeText string

	if getIconVariantFromConfig() == Ascii {
		interactiveTitleText = "[ Igitt Interactive ]\n"
		interactiveByeText = "[ Bye ]\n"
	} else {
		interactiveTitleText = "⌜ Igitt Interactive ⌟\n"
		interactiveByeText = "⌞ Bye ⌝\n"
	}

	interactiveTitle(interactiveTitleText)
	var allCommands []Command
	var commands []Command

	formGroups := make(map[string]*huh.Form)

	err := json.Unmarshal(commandJSON, &allCommands)

	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	for _, command := range allCommands {
		if getShowAllCommandsFromConfig() {
			commands = allCommands
			break
		}

		if command.InsideRepoOnly && !utilities.CheckIsRepo() {
			continue
		}

		if command.OutsideRepoOnly && utilities.CheckIsRepo() {
			continue
		}

		commands = append(commands, command)
	}

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
					Description("\n  Select a branch\n").
					OptionsFunc(getBranchOptions, &commandFlowResult.SelectedCommand).
					Value(&commandFlowResult.SelectedBranch))).WithTheme(theme)

	formGroups["ns-choose-branch-action"] =
		huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Branch action selection").
					Description(fmt.Sprintf("\n  Select an action for the branch %s\n", commandFlowResult.SelectedBranch)).
					Options(getBranchActionOptions()...).
					Value(&commandFlowResult.BranchAction)),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Delete branch").
					Description(fmt.Sprintf("\n  Are you sure you want to delete the branch %s?\n", commandFlowResult.SelectedBranch)).
					Value(&commandFlowResult.DeleteBranchConfirm),
			).WithHideFunc(func() bool {
				return commandFlowResult.DeleteBranchConfirm ||
					commandFlowResult.SelectedCommand.NextStep != "ns-choose-branch-action" ||
					commandFlowResult.BranchAction != "Delete"
			})).WithTheme(theme)

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
					Description("\n" + icons.GetLinkIcon(getIconVariantFromConfig()) + "Enter the link to your repository here.\n").
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
					Description("\n" +
						icons.GetCommitIcon(getIconVariantFromConfig()) +
						"Type a short description to the commit.\n\n" +
						icons.GetBranchIcon(getIconVariantFromConfig()) + "  " + git.GetBranches().CheckedOutBranch + "\n" +
						"Files changed: " + git.GetStagedModificationCountAsString() + "\n").
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
							icons.GetNextStepIcon(getIconVariantFromConfig()),
							"Next step: "+commandFlowResult.SelectedCommand.NextStepTitle,
						)
					}

					return fmt.Sprintf(
						"\n%s: %s\n\n   %s  %s next steps\n\n",
						getTitle(commandFlowResult.SelectedCommand),
						commandFlowResult.SelectedCommand.Description,
						icons.GetNoNextStepIcon(getIconVariantFromConfig()),
						"No",
					)
				}, &commandFlowResult.SelectedCommand).
				Options(commandOptions...).
				Value(&commandFlowResult.SelectedCommand),
		),
	).WithTheme(theme).WithHeight(len(commands) + 9)

	err = mainForm.Run()
	if err != nil {
		fmt.Println(color.HiRedString(interactiveByeText))
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

		if commandFlowResult.SelectedBranch == "[newBranch]" {
			commandFlowResult.SelectedCommand.NextStep = "ns-enter-new-branch-name"
			err = runNextStep(formGroups)
			return err
		}

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

		if commandFlowResult.BranchAction == "Delete" && !isCheckedOutAlready {
			logger.InfoLogger.Println("delete command selected, sending to operations")
			git.DeleteBranch(checkedOutBranchWithoutStar)
			return
		}

		if commandFlowResult.BranchAction == "Delete" && isCheckedOutAlready {
			logger.InfoLogger.Println("delete command selected, not sending to operations, branch already checked out")
			fmt.Println("Cannot delete the branch you are currently on")
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

	if commandFlowResult.SelectedCommand.Id == "op-add" && commandFlowResult.GitAddArguments == "" {
		logger.InfoLogger.Println("add command selected with no arguments, sending to operations")
		git.AddEverything()
		return
	}

	if commandFlowResult.SelectedCommand.Id == "op-add" && commandFlowResult.GitAddArguments != "" {
		logger.InfoLogger.Println("add command selected, sending to operations")
		git.AddChanges(commandFlowResult.GitAddArguments)
		return
	}

	if commandFlowResult.SelectedCommand.Id == "igitt-config" {
		logger.InfoLogger.Println("config command selected, printing config path")
		config.GetConfigPath(true)
		return
	}
}
