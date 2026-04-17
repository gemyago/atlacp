package main

import (
	"errors"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newPRCmd(container *dig.Container) *cobra.Command {
	pr := &cobra.Command{
		Use:   "pr",
		Short: "Bitbucket pull request operations",
	}
	pr.AddCommand(newPRReadCmd(container))
	return pr
}

func newPRReadCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get pull request details",
		RunE: func(_ *cobra.Command, _ []string) error {
			return container.Invoke(func(_ *app.BitbucketService) error {
				if noop {
					return nil
				}
				return errors.New("pr read: not yet implemented")
			})
		},
	}
	cmd.Flags().String("repo-owner", "", "Repository owner (workspace)")
	cmd.Flags().String("repo-name", "", "Repository name (slug)")
	cmd.Flags().Int("pr-id", 0, "Pull request ID")
	_ = cmd.MarkFlagRequired("repo-owner")
	_ = cmd.MarkFlagRequired("repo-name")
	_ = cmd.MarkFlagRequired("pr-id")
	return cmd
}
