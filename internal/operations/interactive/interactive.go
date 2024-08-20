package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func StartInteractive(rootCmd *cobra.Command) {
	fmt.Println("Entering interactive mode. Type 'exit' or 'quit' to leave.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting interactive mode.")
			os.Exit(0)
		}

		args := strings.Split(input, " ")
		rootCmd.SetArgs(args)

		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}
