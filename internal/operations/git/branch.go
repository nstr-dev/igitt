package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/noahstreller/igitt/internal/utilities"
	"github.com/noahstreller/igitt/internal/utilities/logger"
)

type BranchResult struct {
	Branches         []string
	CheckedOutBranch string
}

func GetBranches() BranchResult {
	var branches []string

	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "branch", "-l").CombinedOutput()
	progressIndicator.Stop()

	branchesAsString := string(byteOut)
	branchesAsString = utilities.RemoveLastEmptyLine(branchesAsString)

	branches = strings.Split(branchesAsString, "\n")

	checkedOutBranch := GetCheckedOutBranch(branches)

	branchesTrimmed := trimBranchPrefixes(branches)

	if errOut != nil {
		logger.ErrorLogger.Println("Error:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return BranchResult{}
	}

	logger.InfoLogger.Println("Branch:", errOut, string(byteOut))

	return BranchResult{
		Branches:         branchesTrimmed,
		CheckedOutBranch: checkedOutBranch,
	}
}

func GetCheckedOutBranch(branches []string) string {
	for _, branch := range branches {
		if strings.HasPrefix(branch, "*") {
			return strings.TrimSpace(strings.TrimPrefix(branch, "*"))
		}
	}

	return ""
}

func trimBranchPrefix(branch string) string {
	return strings.TrimSpace(strings.TrimPrefix(branch, "*"))
}

func trimBranchPrefixes(branches []string) []string {
	var trimmedBranches []string

	for _, branch := range branches {
		trimmedBranches = append(trimmedBranches, trimBranchPrefix(branch))
	}

	return trimmedBranches
}

func CheckoutBranch(branch string) {
	fmt.Println("Checking out branch:", color.HiGreenString(branch))
	progressIndicator := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	progressIndicator.Start()
	byteOut, errOut := exec.Command("git", "checkout", branch).CombinedOutput()
	progressIndicator.Stop()

	if errOut != nil {
		logger.ErrorLogger.Println("Error checking out:", errOut, string(byteOut))
		utilities.PrintError(string(byteOut))
		return
	}

	logger.InfoLogger.Println("Checkout:", errOut, string(byteOut))
}
