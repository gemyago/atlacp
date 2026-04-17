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

func bindPRCoreFlags(cmd *cobra.Command) {
	cmd.Flags().String("repo-owner", "", "Repository owner (workspace)")
	cmd.Flags().String("repo-name", "", "Repository name (slug)")
	cmd.Flags().Int("pr-id", 0, "Pull request ID")
	cmd.Flags().String("account", "", "Atlassian account name (optional)")
}

func requirePRCoreFlags(cmd *cobra.Command) {
	_ = cmd.MarkFlagRequired("repo-owner")
	_ = cmd.MarkFlagRequired("repo-name")
	_ = cmd.MarkFlagRequired("pr-id")
}

func parsePRCoreFlags(cmd *cobra.Command) (prCoreIDs, error) {
	var z prCoreIDs
	var err error
	z.RepoOwner, err = cmd.Flags().GetString("repo-owner")
	if err != nil {
		return z, err
	}
	z.RepoName, err = cmd.Flags().GetString("repo-name")
	if err != nil {
		return z, err
	}
	z.PRID, err = cmd.Flags().GetInt("pr-id")
	if err != nil {
		return z, err
	}
	z.Account, err = cmd.Flags().GetString("account")
	return z, err
}
