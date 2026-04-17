package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newPRCmd(container *dig.Container) *cobra.Command {
	pr := &cobra.Command{
		Use:   "pr",
		Short: "Bitbucket pull request operations",
	}
	pr.AddCommand(
		newPRCreateCmd(container),
		newPRReadCmd(container),
		newPRUpdateCmd(container),
		newPRApproveCmd(container),
		newPRRequestChangesCmd(container),
		newPRMergeCmd(container),
		newPRListTasksCmd(container),
		newPRCreateTaskCmd(container),
		newPRUpdateTaskCmd(container),
		newPRDiffStatCmd(container),
		newPRDiffCmd(container),
		newPRAddCommentCmd(container),
		newPRListCommentsCmd(container),
		newPRResolveCommentCmd(container),
	)
	return pr
}

func writeJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	if _, werr := fmt.Fprintln(cmd.OutOrStdout(), string(data)); werr != nil {
		return fmt.Errorf("write JSON: %w", werr)
	}
	return nil
}

type prCreateOpts struct {
	Title          string
	SourceBranch   string
	TargetBranch   string
	RepoOwner      string
	RepoName       string
	Description    string
	Account        string
}

func newPRCreateCmd(container *dig.Container) *cobra.Command {
	var opts prCreateOpts
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRCreate(cmd, container, app.BitbucketCreatePRParams{
				Title:        opts.Title,
				SourceBranch: opts.SourceBranch,
				DestBranch:   opts.TargetBranch,
				RepoOwner:    opts.RepoOwner,
				RepoName:     opts.RepoName,
				Description:  opts.Description,
				AccountName:  opts.Account,
			})
		},
	}
	cmd.Flags().StringVar(&opts.Title, "title", "", "Pull request title")
	cmd.Flags().StringVar(&opts.SourceBranch, "source-branch", "", "Source branch name")
	cmd.Flags().StringVar(&opts.TargetBranch, "target-branch", "", "Target branch name")
	cmd.Flags().StringVar(&opts.RepoOwner, "repo-owner", "", "Repository owner (workspace)")
	cmd.Flags().StringVar(&opts.RepoName, "repo-name", "", "Repository name (slug)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Pull request description")
	cmd.Flags().StringVar(&opts.Account, "account", "", "Atlassian account name (optional)")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("source-branch")
	_ = cmd.MarkFlagRequired("target-branch")
	_ = cmd.MarkFlagRequired("repo-owner")
	_ = cmd.MarkFlagRequired("repo-name")
	return cmd
}

func runPRCreate(cmd *cobra.Command, container *dig.Container, params app.BitbucketCreatePRParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.CreatePR(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRReadCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get pull request details",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRRead(cmd, container, app.BitbucketReadPRParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				AccountName:   core.Account,
			})
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	return cmd
}

func runPRRead(cmd *cobra.Command, container *dig.Container, params app.BitbucketReadPRParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.ReadPR(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRUpdateCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var upd struct {
		Title       string
		Description string
		Draft       bool
	}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update pull request title, description, or draft state",
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := app.BitbucketUpdatePRParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				AccountName:   core.Account,
			}
			if cmd.Flags().Changed("title") {
				params.Title = upd.Title
			}
			if cmd.Flags().Changed("description") {
				params.Description = upd.Description
			}
			if cmd.Flags().Changed("draft") {
				d := upd.Draft
				params.Draft = &d
			}
			return runPRUpdate(cmd, container, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&upd.Title, "title", "", "Updated title")
	cmd.Flags().StringVar(&upd.Description, "description", "", "Updated description")
	cmd.Flags().BoolVar(&upd.Draft, "draft", false, "Set draft state")
	return cmd
}

func runPRUpdate(cmd *cobra.Command, container *dig.Container, params app.BitbucketUpdatePRParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.UpdatePR(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRApproveCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRApprove(cmd, container, app.BitbucketApprovePRParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				AccountName:   core.Account,
			})
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	return cmd
}

func runPRApprove(cmd *cobra.Command, container *dig.Container, params app.BitbucketApprovePRParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.ApprovePR(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRRequestChangesCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "request-changes",
		Short: "Remove approval / request changes on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRRequestChanges(
				cmd,
				container,
				app.BitbucketRequestPRChangesParams{
					RepoOwner:     core.RepoOwner,
					RepoName:      core.RepoName,
					PullRequestID: core.PRID,
					AccountName:   core.Account,
				},
			)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	return cmd
}

func runPRRequestChanges(
	cmd *cobra.Command,
	container *dig.Container,
	params app.BitbucketRequestPRChangesParams,
) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		status, updatedOn, err := svc.RequestPRChanges(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, struct {
			Status    string    `json:"status"`
			UpdatedOn time.Time `json:"updated_on"`
		}{Status: status, UpdatedOn: updatedOn})
	})
}

func newPRMergeCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var strategy string
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRMerge(cmd, container, app.BitbucketMergePRParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				MergeStrategy: strategy,
				AccountName:   core.Account,
			})
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&strategy, "strategy", "", "Merge strategy: merge_commit, squash, or fast_forward")
	return cmd
}

func runPRMerge(cmd *cobra.Command, container *dig.Container, params app.BitbucketMergePRParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.MergePR(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRListTasksCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "list-tasks",
		Short: "List tasks on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRListTasks(cmd, container, app.BitbucketListTasksParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				AccountName:   core.Account,
			})
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	return cmd
}

func runPRListTasks(cmd *cobra.Command, container *dig.Container, params app.BitbucketListTasksParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.ListTasks(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRCreateTaskCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var taskCreate struct {
		Content   string
		CommentID int64
	}
	cmd := &cobra.Command{
		Use:   "create-task",
		Short: "Create a task on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := app.BitbucketCreateTaskParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				Content:       taskCreate.Content,
				AccountName:   core.Account,
			}
			if cmd.Flags().Changed("comment-id") {
				params.CommentID = taskCreate.CommentID
			}
			return runPRCreateTask(cmd, container, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&taskCreate.Content, "content", "", "Task content")
	cmd.Flags().Int64Var(&taskCreate.CommentID, "comment-id", 0, "Optional comment ID to associate with the task")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func runPRCreateTask(cmd *cobra.Command, container *dig.Container, params app.BitbucketCreateTaskParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.CreateTask(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRUpdateTaskCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var taskUpd struct {
		TaskID  int
		Content string
		State   string
	}
	cmd := &cobra.Command{
		Use:   "update-task",
		Short: "Update a task on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := app.BitbucketUpdateTaskParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				TaskID:        taskUpd.TaskID,
				AccountName:   core.Account,
			}
			if cmd.Flags().Changed("content") {
				params.Content = taskUpd.Content
			}
			if cmd.Flags().Changed("state") {
				params.State = taskUpd.State
			}
			return runPRUpdateTask(cmd, container, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().IntVar(&taskUpd.TaskID, "task-id", 0, "Task ID")
	cmd.Flags().StringVar(&taskUpd.Content, "content", "", "Updated task content")
	cmd.Flags().StringVar(&taskUpd.State, "state", "", "Task state: RESOLVED or UNRESOLVED")
	_ = cmd.MarkFlagRequired("task-id")
	return cmd
}

func runPRUpdateTask(cmd *cobra.Command, container *dig.Container, params app.BitbucketUpdateTaskParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.UpdateTask(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRDiffStatCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "diffstat",
		Short: "List changed files summary for a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRDiffStat(cmd, container, app.BitbucketGetPRDiffStatParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				AccountName:   core.Account,
			})
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	return cmd
}

func runPRDiffStat(cmd *cobra.Command, container *dig.Container, params app.BitbucketGetPRDiffStatParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.GetPRDiffStat(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}

func newPRDiffCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var diffOpts struct {
		Path    string
		Context int
	}
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Get raw diff text for a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := app.BitbucketGetPRDiffParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				AccountName:   core.Account,
			}
			if diffOpts.Path != "" {
				params.FilePaths = []string{diffOpts.Path}
			}
			if cmd.Flags().Changed("context") {
				ctxLines := diffOpts.Context
				params.ContextLines = &ctxLines
			}
			return runPRDiff(cmd, container, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&diffOpts.Path, "path", "", "Optional file path to scope the diff")
	cmd.Flags().IntVar(&diffOpts.Context, "context", 0, "Optional number of context lines")
	return cmd
}

func runPRDiff(cmd *cobra.Command, container *dig.Container, params app.BitbucketGetPRDiffParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.GetPRDiff(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, struct {
			Diff string `json:"diff"`
		}{Diff: result})
	})
}

func newPRAddCommentCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var cmt struct {
		Content  string
		FilePath string
		Line     int
	}
	cmd := &cobra.Command{
		Use:   "add-comment",
		Short: "Post a comment on a pull request (general or inline)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := app.BitbucketAddPRCommentParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				Content:       cmt.Content,
				FilePath:      cmt.FilePath,
				AccountName:   core.Account,
			}
			if cmd.Flags().Changed("line") {
				params.LineFrom = cmt.Line
				params.LineTo = cmt.Line
			}
			return runPRAddComment(cmd, container, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&cmt.Content, "content", "", "Comment content (raw)")
	cmd.Flags().StringVar(&cmt.FilePath, "file-path", "", "File path for inline comments")
	cmd.Flags().IntVar(&cmt.Line, "line", 0, "Line number for inline comments (sets from/to)")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func runPRAddComment(cmd *cobra.Command, container *dig.Container, params app.BitbucketAddPRCommentParams) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		commentID, text, err := svc.AddPRComment(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, struct {
			CommentID int64  `json:"comment_id"`
			Content   string `json:"content"`
		}{CommentID: commentID, Content: text})
	})
}

func newPRListCommentsCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var includeResolved bool
	cmd := &cobra.Command{
		Use:   "list-comments",
		Short: "List comments on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRListComments(
				cmd,
				container,
				app.BitbucketListPRCommentsParams{
					RepoOwner:     core.RepoOwner,
					RepoName:      core.RepoName,
					PullRequestID: core.PRID,
					AccountName:   core.Account,
				},
				includeResolved,
			)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().BoolVar(&includeResolved, "include-resolved", false, "Include resolved comments in the output")
	return cmd
}

func runPRListComments(
	cmd *cobra.Command,
	container *dig.Container,
	params app.BitbucketListPRCommentsParams,
	includeResolved bool,
) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.ListPRComments(cmd.Context(), params)
		if err != nil {
			return err
		}
		if !includeResolved && result != nil {
			filtered := make([]app.BitbucketPRComment, 0, len(result.Values))
			for _, c := range result.Values {
				if !c.Resolved {
					filtered = append(filtered, c)
				}
			}
			result.Values = filtered
		}
		return writeJSON(cmd, result)
	})
}

func newPRResolveCommentCmd(container *dig.Container) *cobra.Command {
	var core prCoreIDs
	var commentID int64
	cmd := &cobra.Command{
		Use:   "resolve-comment",
		Short: "Resolve a pull request comment thread",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRResolveComment(cmd, container, app.BitbucketResolvePRCommentParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				CommentID:     commentID,
				AccountName:   core.Account,
			})
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().Int64Var(&commentID, "comment-id", 0, "Comment ID")
	_ = cmd.MarkFlagRequired("comment-id")
	return cmd
}

func runPRResolveComment(
	cmd *cobra.Command,
	container *dig.Container,
	params app.BitbucketResolvePRCommentParams,
) error {
	return container.Invoke(func(svc *app.BitbucketService) error {
		if noop {
			return nil
		}
		result, err := svc.ResolvePRComment(cmd.Context(), params)
		if err != nil {
			return err
		}
		return writeJSON(cmd, result)
	})
}
