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

type fileContentOpts struct {
	RepoOwner string
	RepoName  string
	Commit    string
	Path      string
	Account   string
}

func newFileContentCmd(container *dig.Container) *cobra.Command {
	var opts fileContentOpts
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Get file content at a commit",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFileContent(cmd, container, app.BitbucketGetFileContentParams{
				AccountName: opts.Account,
				RepoOwner:   opts.RepoOwner,
				RepoName:    opts.RepoName,
				Commit:      opts.Commit,
				Path:        opts.Path,
			})
		},
	}
	cmd.Flags().StringVar(&opts.RepoOwner, "repo-owner", "", "Repository owner (workspace)")
	cmd.Flags().StringVar(&opts.RepoName, "repo-name", "", "Repository name (slug)")
	cmd.Flags().StringVar(&opts.Commit, "commit", "", "Commit hash")
	cmd.Flags().StringVar(&opts.Path, "path", "", "File path in the repository")
	cmd.Flags().StringVar(&opts.Account, "account", "", "Atlassian account name (optional)")
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
