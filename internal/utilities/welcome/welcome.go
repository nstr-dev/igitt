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

	Spacing(1)
	titleMessage("Welcome to Igitt - Interactive Git in the Terminal!")
	Spacing(2)

	heading("Getting Started")
	Spacing(2)
	fmt.Printf("To enter the interactive experience, run %s without any arguments.", command("igitt"))
	Spacing(1)
	fmt.Printf("For a list of available commands, run %s.", command("igitt --help"))
	Spacing(2)

	heading("Alias")
	Spacing(2)
	fmt.Printf("To use the alias %s instead of %s, run following command: %s", command("igt"), command("igitt"), underline(command("igitt igt")))
	Spacing(2)

	heading("Configuration")
	Spacing(2)
	fmt.Printf("To configure Igitt (e.g. to change how icons are displayed and which commands are shown), run %s.", command("igitt config"))
	Spacing(3)
}
