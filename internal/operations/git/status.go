package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

type FileStatus struct {
	StatusLetter string
	FileName     string
}

type ModifiedStatusInfo struct {
	StatusColor  color.Attribute
	StatusTitle  string
	StatusLetter string
}

// https://git-scm.com/docs/git-status#_short_format

var FileStatuses = []ModifiedStatusInfo{
	// Basic States
	{color.FgHiRed, "Deleted (staged)", "D "},
	{color.FgHiGreen, "New (staged)", "A "},
	{color.FgHiYellow, "Modified (staged)", "M "},
	{color.FgHiBlue, "Renamed (staged)", "R "},
	{color.FgHiCyan, "Copied (staged)", "C "},
	{color.FgHiMagenta, "Type changed (staged)", "T "},

	// Work Tree (Unstaged) States
	{color.FgHiYellow, "Modified (unstaged)", " M"},
	{color.FgHiMagenta, "Type changed (unstaged)", " T"},
	{color.FgHiRed, "Deleted (unstaged)", " D"},
	{color.FgHiBlue, "Renamed (unstaged)", " R"},

	// Combined Staged and Unstaged States
	{color.FgHiYellow, "Modified (staged & partially unstaged)", "MM"},
	{color.FgHiMagenta, "Type changed (staged & unstaged)", "TT"},
	{color.FgHiYellow, "New (staged), Modified (unstaged)", "AM"},
	{color.FgHiYellow, "Renamed (staged), Modified (unstaged)", "RM"},
	{color.FgHiYellow, "Copied (staged), Modified (unstaged)", "CM"},
	{color.FgHiMagenta, "Modified (staged), Type changed (unstaged)", "MT"},
	{color.FgHiRed, "Modified (staged), Deleted (unstaged)", "MD"},
	{color.FgHiYellow, "Modified (staged), Renamed (unstaged)", "MR"},
	{color.FgHiYellow, "Modified (staged), Copied (unstaged)", "MC"},
	{color.FgHiRed, "Type changed (staged), Deleted (unstaged)", "TD"},
	{color.FgHiMagenta, "Type changed (staged), Renamed (unstaged)", "TR"},
	{color.FgHiMagenta, "Type changed (staged), Copied (unstaged)", "TC"},
	{color.FgHiMagenta, "Renamed (staged), Type changed (unstaged)", "RT"},
	{color.FgHiRed, "Renamed (staged), Deleted (unstaged)", "RD"},
	{color.FgHiRed, "Copied (staged), Deleted (unstaged)", "CD"},
	{color.FgHiCyan, "Copied (staged), Renamed (unstaged)", "CR"},

	// Unmerged States
	{color.FgHiRed, "Unmerged, Both Deleted", "DD"},
	{color.FgHiGreen, "Unmerged, Added by Us", "AU"},
	{color.FgHiYellow, "Unmerged, Deleted by Them", "UD"},
	{color.FgHiGreen, "Unmerged, Added by Them", "UA"},
	{color.FgHiYellow, "Unmerged, Deleted by Us", "DU"},
	{color.FgHiGreen, "Unmerged, Both Added", "AA"},
	{color.FgHiYellow, "Unmerged, Both Modified", "UU"},

	// Other States
	{color.FgHiBlack, "Untracked", "??"},
	{color.FgHiBlack, "Ignored", "!!"},
}

func Status() {
	modifications, err := getModifications()
	if err != nil {
		logger.ErrorLogger.Println("Failed to get modifications: ", err)
		return
	}

	if len(modifications) == 0 {
		fmt.Println(color.HiGreenString("✓"), "Up to date.")
		return
	}

	fmt.Printf("\nFiles with changes:\n")
	fmt.Printf("===================\n\n")

	statusMap := make(map[string]ModifiedStatusInfo)
	for _, status := range FileStatuses {
		statusMap[status.StatusLetter] = status
	}

	maxWidth := 0

	for _, modification := range modifications {
		if status, exists := statusMap[modification.StatusLetter]; exists {
			color := color.New(status.StatusColor).SprintFunc()
			coloredTitle := color(status.StatusTitle)
			if len(coloredTitle) > maxWidth {
				maxWidth = len(coloredTitle) + 4
			}
		}
	}

	for _, modification := range modifications {
		if status, exists := statusMap[modification.StatusLetter]; exists {
			color := color.New(status.StatusColor).SprintFunc()
			coloredTitle := color(status.StatusTitle)
			format := fmt.Sprintf("%%-%ds%%s\n", maxWidth)
			fmt.Printf(format, coloredTitle, modification.FileName)
		}
	}
}

func getModifications() ([]FileStatus, error) {
	var statuses []FileStatus

	status, err := runGitStatus()
	if err != nil {
		logger.ErrorLogger.Println("Failed previous step, aborting: ", err)
		utilities.PrintError(status)
		return nil, err
	}

	statusWithoutEmptyLine := utilities.RemoveLastEmptyLine(status)
	statusLines := strings.Split(statusWithoutEmptyLine, "\n")

	if len(statusLines) == 0 {
		return statuses, nil
	}

	if len(statusLines) == 1 {
		if statusLines[0] == "" {
			return []FileStatus{}, nil
		}
	}

	for _, line := range statusLines {
		statusLetter := string(line[0:2])
		fileName := line[3:]

		statuses = append(statuses, FileStatus{statusLetter, fileName})
	}

	statusTitleMap := make(map[string]string)
	for _, status := range FileStatuses {
		statusTitleMap[status.StatusLetter] = status.StatusTitle
	}

	sort.Slice(statuses, func(i, j int) bool {
		return statusTitleMap[statuses[i].StatusLetter] < statusTitleMap[statuses[j].StatusLetter]
	})

	return statuses, nil
}

func runGitStatus() (string, error) {
	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "status", "--porcelain").CombinedOutput()
	progressIndicator.Stop()

	logger.InfoLogger.Println("Fetching git status:", errOut, string(byteOut))

	return string(byteOut), errOut
}
