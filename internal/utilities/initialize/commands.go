package initialize

import (
	"os"

	"github.com/noahstreller/igitt/internal/operations/git"
	"github.com/noahstreller/igitt/internal/operations/interactive"
	"github.com/noahstreller/igitt/internal/utilities/logger"
	"github.com/spf13/cobra"
)

func InitializeIgitt() {
	var rootCmd = &cobra.Command{
		Use:   "igitt",
		Short: "Igitt is a supercharged Git client with a CLI.",
		Long:  `Igitt supercharges your Git experience with an interactive CLI. Designed to enhance learning and streamline workflows, it offers detailed command descriptions and efficient shortcuts for a faster, more intuitive Git journey.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.InfoLogger.Println("igitt was called without arguments")
		},
	}

	var cloneCmd = &cobra.Command{
		Use:     "clone [repository]",
		Short:   "(cln) Clone a repository into a new directory",
		Aliases: []string{"cln"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			git.CloneRepository(args[0])
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create an empty Git repository or reinitialize an existing one",
		Run: func(cmd *cobra.Command, args []string) {
			git.InitRepository()
		},
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Fetch from and integrate with another repository or a local branch",
		Run: func(cmd *cobra.Command, args []string) {
			git.PullRemote()
		},
	}

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Update remote refs along with associated objects",
		Run: func(cmd *cobra.Command, args []string) {
			git.PushRemote()
		},
	}

	var interactiveCmd = &cobra.Command{
		Use:     "interactive",
		Short:   "(i) Enter interactive mode",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			interactive.StartInteractive(rootCmd)
		},
	}

	cobra.OnInitialize(func() {
		if len(os.Args) == 1 {
			rootCmd.Help()
			os.Exit(0)
		}
	})

	rootCmd.AddCommand(
		cloneCmd,
		initCmd,
		interactiveCmd,
		pullCmd,
		pushCmd,
	)
	rootCmd.Execute()
}
