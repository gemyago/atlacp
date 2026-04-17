package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newFileCmd(_ *dig.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "file",
		Short: "Bitbucket file operations",
	}
}
