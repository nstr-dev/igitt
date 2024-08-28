package initialize

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/noahstreller/igitt/internal/operations"
	"github.com/noahstreller/igitt/internal/operations/git"
	"github.com/noahstreller/igitt/internal/operations/interactive"
	"github.com/noahstreller/igitt/internal/utilities/logger"
	"github.com/spf13/cobra"
)

func InitializeIgitt() {
	var rootCmd = &cobra.Command{
		Use:   "igitt",
		Short: "Igitt is an interactive Git client with a CLI.",
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

	var commitCmd = &cobra.Command{
		Use:     "commit",
		Short:   "(cmt) Record changes to the repository",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"cmt"},
		Run: func(cmd *cobra.Command, args []string) {
			git.CommitChanges(args[0])
		},
	}

	var gitAddCmd = &cobra.Command{
		Use:     "add",
		Short:   "(a, +) Add file contents to the index",
		Aliases: []string{"a", "+"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			git.AddChanges(args[0])
		},
	}

	var statusCmd = &cobra.Command{
		Use:     "status",
		Short:   "(s) Show changed files",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			git.Status()
		},
	}

	var interactiveCmd = &cobra.Command{
		Use:     "interactive",
		Short:   "(i) Enter interactive mode",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			interactive.StartInteractive()
		},
	}

	var createAliasScripts = &cobra.Command{
		Use:     "mkalias",
		Short:   "Create scripts to use igt as an alias for igitt",
		Aliases: []string{"igt"},
		Run: func(cmd *cobra.Command, args []string) {
			operations.CreateAliasScripts()
		},
	}

	cobra.OnInitialize(func() {
		if len(os.Args) == 1 {
			fmt.Println(color.RedString("No arguments provided."))
			// interactive.StartInteractive()
		}
	})

	rootCmd.AddCommand(
		cloneCmd,
		initCmd,
		interactiveCmd,
		gitAddCmd,
		pullCmd,
		pushCmd,
		statusCmd,
		commitCmd,
		createAliasScripts,
	)
	err := rootCmd.Execute()
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
}
