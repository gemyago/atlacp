package main

import (
	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newFileCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	file := &cobra.Command{
		Use:     "file",
		Short:   "Bitbucket file operations",
		Long:    "Read file-level data from Bitbucket repositories and commits.",
		Example: `bbmd file content --repo-owner <workspace> --repo-name <repo> --commit <sha> --path <file-path>`,
	}
	file.AddCommand(newFileContentCmd(container, rootParams))
	return file
}

type fileContentOpts struct {
	RepoOwner string
	RepoName  string
	Commit    string
	Path      string
	Account   string
}

func newFileContentCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var opts fileContentOpts
	cmd := &cobra.Command{
		Use:     "content",
		Short:   "Get file content at a commit",
		Long:    "Read file content from a repository at a specific commit hash.",
		Example: `bbmd file content --repo-owner <workspace> --repo-name <repo> --commit <sha> --path <file path>`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFileContent(cmd, container, rootParams, app.BitbucketGetFileContentParams{
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

func runFileContent(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketGetFileContentParams,
) error {
	return container.Invoke(func(execDeps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(
			cmd,
			execDeps,
			execArgs[app.BitbucketGetFileContentParams, *bitbucket.FileContentResult]{
				rootParams: rootParams,
				params:     params,
				target:     svc.GetFileContent,
			})
	})
}
