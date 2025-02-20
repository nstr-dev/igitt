package welcome

import (
	"fmt"

	"github.com/fatih/color"
)

func Spacing(amount int) {
	for i := 0; i < amount; i++ {
		fmt.Println()
	}
}

func PrintWelcomeMessage() {
	titleMessage := color.New(color.Bold, color.FgGreen).PrintfFunc()
	heading := color.New(color.Italic, color.Underline).PrintfFunc()
	command := color.New(color.FgCyan, color.Bold).SprintfFunc()
	underline := color.New(color.Underline).SprintfFunc()
	titleMessage("Welcome to Igitt - Interactive Git in the Terminal!")
	Spacing(2)

	heading("Getting Started")
	Spacing(2)
	fmt.Println("To enter the interactive experience, run igitt without any arguments.")
	fmt.Println("For a list of available commands, run igitt --help.")
	Spacing(1)

	heading("Alias")
	Spacing(2)
	fmt.Printf("To use the alias %s instead of %s, run following command: %s", command("igt"), command("igitt"), underline(command("igitt igt")))
}
