package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/syst3mctl/go-ctl/internal/cli/completion"
	"github.com/syst3mctl/go-ctl/internal/cli/help"
)

// NewCompletionCommand creates the completion command with enhanced features
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script with dynamic suggestions",
		Long: `Generate the autocompletion script for the specified shell with enhanced
dynamic completion support.

This command generates a completion script that enables intelligent tab completion
for go-ctl commands, flags, and dynamic values like package names, template names,
and configuration options.

FEATURES:
• Context-aware completion based on command and flags
• Dynamic package name suggestions from popular Go packages
• Template name completion with descriptions
• Configuration value completion (Go versions, frameworks, etc.)
• File and directory path completion where appropriate

INSTALLATION:

# Bash
go-ctl completion bash > /etc/bash_completion.d/go-ctl
# Or for user-specific installation:
go-ctl completion bash > ~/.bash_completions/go-ctl

# Zsh
go-ctl completion zsh > "${fpath[1]}/_go-ctl"
# Or add to your .zshrc:
echo 'source <(go-ctl completion zsh)' >> ~/.zshrc

# Fish
go-ctl completion fish > ~/.config/fish/completions/go-ctl.fish

# PowerShell
go-ctl completion powershell > go-ctl.ps1
# Then source it in your PowerShell profile

EXAMPLES:
  # Generate bash completion with enhanced features
  go-ctl completion bash

  # Test completion (after installation)
  go-ctl generate my-api --http=<TAB>  # Shows: gin, echo, fiber, chi, net-http
  go-ctl template show <TAB>           # Shows: minimal, api, microservice, cli
  go-ctl package search <TAB>          # Shows popular package suggestions

For more examples and setup instructions:
https://github.com/syst3mctl/go-ctl/blob/main/docs/shell-completion.md`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:                  runCompletion,
	}

	// Add enhanced help
	help.AddEnhancedHelp(cmd)

	// Setup dynamic completion for the root command
	rootCmd := cmd.Root()
	completion.SetupDynamicCompletion(rootCmd)

	return cmd
}

// runCompletion generates shell completion scripts with enhanced features
func runCompletion(cmd *cobra.Command, args []string) error {
	printInfo("Generating %s completion script with enhanced features...", args[0])

	var err error
	switch args[0] {
	case "bash":
		err = cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		err = cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		err = cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return cmd.Help()
	}

	if err != nil {
		return err
	}

	if !isQuiet() {
		printSuccess("Completion script generated successfully!")
		printInfo("Enhanced features include:")
		printInfo("  • Dynamic package name suggestions")
		printInfo("  • Template completion with descriptions")
		printInfo("  • Context-aware flag completion")
		printInfo("  • File and directory path completion")
		printInfo("")
		printInfo("After installation, try these examples:")
		printInfo("  go-ctl generate --http=<TAB>")
		printInfo("  go-ctl template show <TAB>")
		printInfo("  go-ctl package search <TAB>")
	}

	return nil
}
