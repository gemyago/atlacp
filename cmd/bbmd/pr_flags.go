package main

import (
	"github.com/spf13/cobra"
)

// prCoreIDs holds standard repo + PR flags shared by most `bbmd pr` sub-commands.
type prCoreIDs struct {
	RepoOwner string
	RepoName  string
	PRID      int
	Account   string
}

func bindPRCoreFlags(cmd *cobra.Command, core *prCoreIDs) {
	cmd.Flags().StringVar(&core.RepoOwner, "repo-owner", "", "Repository owner (workspace)")
	cmd.Flags().StringVar(&core.RepoName, "repo-name", "", "Repository name (slug)")
	cmd.Flags().IntVar(&core.PRID, "pr-id", 0, "Pull request ID")
	cmd.Flags().StringVar(&core.Account, "account", "", "Atlassian account name (optional)")
}

func requirePRCoreFlags(cmd *cobra.Command) {
	_ = cmd.MarkFlagRequired("repo-owner")
	_ = cmd.MarkFlagRequired("repo-name")
	_ = cmd.MarkFlagRequired("pr-id")
}
