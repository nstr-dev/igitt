package initialize

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nstr-dev/igitt/internal/operations"
	"github.com/nstr-dev/igitt/internal/operations/git"
	"github.com/nstr-dev/igitt/internal/operations/interactive"
	"github.com/nstr-dev/igitt/internal/utilities/config"
	"github.com/nstr-dev/igitt/internal/utilities/logger"
	"github.com/nstr-dev/igitt/internal/utilities/welcome"
	"github.com/spf13/cobra"
)

func InitializeIgitt(version string, commit string, buildDate string) {
	cyan := color.New(color.FgCyan).SprintfFunc()
	heading := color.New(color.Bold, color.FgGreen).SprintfFunc()

	var rootCmd = &cobra.Command{
		Use:   "igitt",
		Short: "Igitt is an interactive Git client with a CLI.",
		Long:  `Igitt supercharges your Git experience with an interactive CLI. Designed to enhance learning and streamline workflows, it offers detailed command descriptions and efficient shortcuts for a faster, more intuitive Git journey.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.InfoLogger.Println("igitt was called without arguments")
		},
		Version: ">> Version: " + cyan(version) + "\n>> Commit: " + cyan(commit) + "\n>> Build Date: " + cyan(buildDate),
	}

	rootCmd.SetVersionTemplate(
		heading("Igitt - Interactive Git in the Terminal") +
			"\n=======================================\n\n{{.Version}}",
	)

	var cloneCmd = &cobra.Command{
		Use:     "clone [repository]",
		Short:   "(cln) Clone a repository into a new directory",
		Aliases: []string{"cln"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			git.CloneRepository(strings.Join(args, " "))
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
			git.CommitChanges(strings.Join(args, " "))
		},
	}

	var gitAddCmd = &cobra.Command{
		Use:     "add",
		Short:   "(a, +) Add file contents to the index",
		Aliases: []string{"a", "+"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			git.AddChanges(args)
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

	var checkoutCmd = &cobra.Command{
		Use:     "checkout",
		Short:   "(cout) Change to a different branch",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"cout"},
		Run: func(cmd *cobra.Command, args []string) {
			git.CheckoutBranch(strings.Join(args, " "))
		},
	}

	var branchCmd = &cobra.Command{
		Use:     "branch",
		Short:   "(br) Manage branches",
		Aliases: []string{"br"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				git.DoCustomBranchAction("")
				return
			}
			git.DoCustomBranchAction(strings.Join(args, " "))
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

	var igittConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Print the path of the igitt configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			config.GetConfigPath(true)
		},
	}

	cobra.OnInitialize(func() {
		createdNewConfig, err := config.InitialConfig()

		if err != nil {
			logger.ErrorLogger.Fatal(err)
			return
		}

		if len(os.Args) == 1 {
			if createdNewConfig {
				welcome.PrintWelcomeMessage()
				return
			}
			interactive.StartInteractive()
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
		checkoutCmd,
		commitCmd,
		createAliasScripts,
		branchCmd,
		igittConfigCmd,
	)
	err := rootCmd.Execute()
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
}
