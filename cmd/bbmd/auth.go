package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newAuthCmd(_ *dig.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Atlassian account management",
	}
}
