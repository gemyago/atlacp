package main

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func setupCommands() *cobra.Command {
	container := dig.New()
	return newRootCmd(container, new(rootCommandParams))
}

func main() { // coverage-ignore
	rootCmd := setupCommands()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
