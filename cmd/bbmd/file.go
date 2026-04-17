package main

import (
	"github.com/gemyago/atlacp/internal/app"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newFileCmd(container *dig.Container) *cobra.Command {
	file := &cobra.Command{
		Use:   "file",
		Short: "Bitbucket file operations",
	}
	file.AddCommand(newFileContentCmd(container))
	return file
}

func newFileContentCmd(container *dig.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Get file content at a commit",
		RunE: func(cmd *cobra.Command, _ []string) error {
			repoOwner, _ := cmd.Flags().GetString("repo-owner")
			repoName, _ := cmd.Flags().GetString("repo-name")
			commit, _ := cmd.Flags().GetString("commit")
			path, _ := cmd.Flags().GetString("path")
			account, _ := cmd.Flags().GetString("account")
			return runFileContent(cmd, container, app.BitbucketGetFileContentParams{
				AccountName: account,
				RepoOwner:   repoOwner,
				RepoName:    repoName,
				Commit:      commit,
				Path:        path,
			})
		},
	}
	cmd.Flags().String("repo-owner", "", "Repository owner (workspace)")
	cmd.Flags().String("repo-name", "", "Repository name (slug)")
	cmd.Flags().String("commit", "", "Commit hash")
	cmd.Flags().String("path", "", "File path in the repository")
	cmd.Flags().String("account", "", "Atlassian account name (optional)")
	_ = cmd.MarkFlagRequired("repo-owner")
	_ = cmd.MarkFlagRequired("repo-name")
	_ = cmd.MarkFlagRequired("commit")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func runFileContent(cmd *cobra.Command, container *dig.Container, params app.BitbucketGetFileContentParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.GetFileContent(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}
