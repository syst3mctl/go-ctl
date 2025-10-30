package main

import (
	"fmt"
	"os"

	"github.com/syst3mctl/go-ctl/internal/cli/commands"
)

func main() {
	rootCmd := commands.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
