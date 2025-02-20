package welcome

import "fmt"

func WelcomeMessage() string {
	return "Welcome to Igitt! This is a CLI tool to manage your git repositories. For a list of commands, type 'igitt help'."
}

func PrintWelcomeMessage() {
	fmt.Println(WelcomeMessage())
}
