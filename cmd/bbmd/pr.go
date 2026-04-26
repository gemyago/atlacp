package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type prRequestChangesOut struct {
	Status    string    `json:"status"`
	UpdatedOn time.Time `json:"updated_on"`
}

type prDiffOut struct {
	Diff string `json:"diff"`
}

type prAddCommentOut struct {
	CommentID int64  `json:"comment_id"`
	Content   string `json:"content"`
}

type prListCommentsCLI struct {
	Params          app.BitbucketListPRCommentsParams
	IncludeResolved bool
}

// applyPRCommentsResolvedFilter drops resolved comments when includeResolved is false.
func applyPRCommentsResolvedFilter(result *app.BitbucketListPRCommentsResult, includeResolved bool) {
	if includeResolved || result == nil {
		return
	}
	filtered := make([]app.BitbucketPRComment, 0, len(result.Values))
	for _, c := range result.Values {
		if !c.Resolved {
			filtered = append(filtered, c)
		}
	}
	result.Values = filtered
}

func newPRCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	pr := &cobra.Command{
		Use:   "pr",
		Short: "Bitbucket pull request operations",
		Long:  "Manage Bitbucket pull requests: create, read, update, merge, and comment workflows.",
		Example: `bbmd pr list-tasks --repo-owner <workspace> --repo-name <repo> --pr-id <id>
bbmd pr read --repo-owner <workspace> --repo-name <repo> --pr-id <id>`,
	}
	pr.AddCommand(
		newPRCreateCmd(container, rootParams),
		newPRReadCmd(container, rootParams),
		newPRUpdateCmd(container, rootParams),
		newPRApproveCmd(container, rootParams),
		newPRRequestChangesCmd(container, rootParams),
		newPRMergeCmd(container, rootParams),
		newPRListTasksCmd(container, rootParams),
		newPRCreateTaskCmd(container, rootParams),
		newPRUpdateTaskCmd(container, rootParams),
		newPRDiffStatCmd(container, rootParams),
		newPRDiffCmd(container, rootParams),
		newPRAddCommentCmd(container, rootParams),
		newPRListCommentsCmd(container, rootParams),
		newPRResolveCommentCmd(container, rootParams),
	)
	return pr
}

type prCreateOpts struct {
	Title        string
	SourceBranch string
	TargetBranch string
	RepoOwner    string
	RepoName     string
	Description  string
	Account      string
}

func newPRCreateCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var opts prCreateOpts
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a pull request",
		Long:    "Create a new pull request from a source branch to a target branch.",
		Example: `bbmd pr create --repo-owner <workspace> --repo-name <repo> --title <title> --source-branch <branch> --target-branch <branch>`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRCreate(cmd, container, rootParams, app.BitbucketCreatePRParams{
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

func runPRCreate(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketCreatePRParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketCreatePRParams, *bitbucket.PullRequest]{
			rootParams: rootParams,
			params:     params,
			target:     svc.CreatePR,
		})
	})
}

func newPRReadCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:     "read",
		Short:   "Get pull request details",
		Long:    "Read detailed metadata for a pull request.",
		Example: `bbmd pr read --repo-owner <workspace> --repo-name <repo> --pr-id <id>`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRRead(cmd, container, rootParams, app.BitbucketReadPRParams{
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

func runPRRead(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketReadPRParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketReadPRParams, *bitbucket.PullRequest]{
			rootParams: rootParams,
			params:     params,
			target:     svc.ReadPR,
		})
	})
}

func newPRUpdateCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
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
			return runPRUpdate(cmd, container, rootParams, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&upd.Title, "title", "", "Updated title")
	cmd.Flags().StringVar(&upd.Description, "description", "", "Updated description")
	cmd.Flags().BoolVar(&upd.Draft, "draft", false, "Set draft state")
	return cmd
}

func runPRUpdate(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketUpdatePRParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketUpdatePRParams, *bitbucket.PullRequest]{
			rootParams: rootParams,
			params:     params,
			target:     svc.UpdatePR,
		})
	})
}

func newPRApproveCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRApprove(cmd, container, rootParams, app.BitbucketApprovePRParams{
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

func runPRApprove(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketApprovePRParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketApprovePRParams, *bitbucket.Participant]{
			rootParams: rootParams,
			params:     params,
			target:     svc.ApprovePR,
		})
	})
}

func newPRRequestChangesCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "request-changes",
		Short: "Remove approval / request changes on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRRequestChanges(
				cmd,
				container,
				rootParams,
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
	rootParams *rootCommandParams,
	params app.BitbucketRequestPRChangesParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketRequestPRChangesParams, prRequestChangesOut]{
			rootParams: rootParams,
			params:     params,
			target: func(ctx context.Context, p app.BitbucketRequestPRChangesParams) (prRequestChangesOut, error) {
				status, updatedOn, err := svc.RequestPRChanges(ctx, p)
				if err != nil {
					return prRequestChangesOut{}, err
				}
				return prRequestChangesOut{Status: status, UpdatedOn: updatedOn}, nil
			},
		})
	})
}

func newPRMergeCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	var strategy string
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRMerge(cmd, container, rootParams, app.BitbucketMergePRParams{
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

func runPRMerge(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketMergePRParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketMergePRParams, *bitbucket.PullRequest]{
			rootParams: rootParams,
			params:     params,
			target:     svc.MergePR,
		})
	})
}

func newPRListTasksCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "list-tasks",
		Short: "List tasks on a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRListTasks(cmd, container, rootParams, app.BitbucketListTasksParams{
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

func runPRListTasks(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketListTasksParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketListTasksParams, *bitbucket.PaginatedTasks]{
			rootParams: rootParams,
			params:     params,
			target:     svc.ListTasks,
		})
	})
}

func newPRCreateTaskCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
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
			return runPRCreateTask(cmd, container, rootParams, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&taskCreate.Content, "content", "", "Task content")
	cmd.Flags().Int64Var(&taskCreate.CommentID, "comment-id", 0, "Optional comment ID to associate with the task")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func runPRCreateTask(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketCreateTaskParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketCreateTaskParams, *bitbucket.PullRequestCommentTask]{
			rootParams: rootParams,
			params:     params,
			target:     svc.CreateTask,
		})
	})
}

func newPRUpdateTaskCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
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
			return runPRUpdateTask(cmd, container, rootParams, params)
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

func runPRUpdateTask(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketUpdateTaskParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketUpdateTaskParams, *bitbucket.PullRequestCommentTask]{
			rootParams: rootParams,
			params:     params,
			target:     svc.UpdateTask,
		})
	})
}

func newPRDiffStatCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	cmd := &cobra.Command{
		Use:   "diffstat",
		Short: "List changed files summary for a pull request",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRDiffStat(cmd, container, rootParams, app.BitbucketGetPRDiffStatParams{
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

func runPRDiffStat(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketGetPRDiffStatParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketGetPRDiffStatParams, *app.PaginatedDiffStat]{
			rootParams: rootParams,
			params:     params,
			target:     svc.GetPRDiffStat,
		})
	})
}

func newPRDiffCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
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
			return runPRDiff(cmd, container, rootParams, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&diffOpts.Path, "path", "", "Optional file path to scope the diff")
	cmd.Flags().IntVar(&diffOpts.Context, "context", 0, "Optional number of context lines")
	return cmd
}

func runPRDiff(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketGetPRDiffParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketGetPRDiffParams, prDiffOut]{
			rootParams: rootParams,
			params:     params,
			target: func(ctx context.Context, p app.BitbucketGetPRDiffParams) (prDiffOut, error) {
				s, err := svc.GetPRDiff(ctx, p)
				if err != nil {
					return prDiffOut{}, err
				}
				return prDiffOut{Diff: s}, nil
			},
		})
	})
}

func newPRAddCommentCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	var cmt struct {
		Content  string
		FilePath string
		Line     int
		ParentID int64
	}
	cmd := &cobra.Command{
		Use:   "add-comment",
		Short: "Post a comment on a pull request (general or inline)",
		Long:  "Post a top-level or inline comment on a pull request.",
		Example: `bbmd pr add-comment --repo-owner <workspace> --repo-name <repo> --pr-id <id> --content "<comment>"
bbmd pr add-comment --repo-owner <workspace> --repo-name <repo> --pr-id <id> --content "<comment>" --file-path <path> --line <line> --parent-comment-id <id>`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			params := app.BitbucketAddPRCommentParams{
				RepoOwner:     core.RepoOwner,
				RepoName:      core.RepoName,
				PullRequestID: core.PRID,
				Content:       cmt.Content,
				FilePath:      cmt.FilePath,
				AccountName:   core.Account,
				ParentID:      cmt.ParentID,
			}
			if cmd.Flags().Changed("line") {
				params.LineFrom = cmt.Line
				params.LineTo = cmt.Line
			}
			return runPRAddComment(cmd, container, rootParams, params)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().StringVar(&cmt.Content, "content", "", "Comment content (raw)")
	cmd.Flags().StringVar(&cmt.FilePath, "file-path", "", "File path for inline comments")
	cmd.Flags().IntVar(&cmt.Line, "line", 0, "Line number for inline comments (sets from/to)")
	cmd.Flags().Int64Var(&cmt.ParentID, "parent-comment-id", 0, "Parent comment ID for replies (optional)")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func runPRAddComment(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	params app.BitbucketAddPRCommentParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketAddPRCommentParams, prAddCommentOut]{
			rootParams: rootParams,
			params:     params,
			target: func(ctx context.Context, p app.BitbucketAddPRCommentParams) (prAddCommentOut, error) {
				id, text, err := svc.AddPRComment(ctx, p)
				if err != nil {
					return prAddCommentOut{}, err
				}
				return prAddCommentOut{CommentID: id, Content: text}, nil
			},
		})
	})
}

func newPRListCommentsCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	var includeResolved bool
	var page, pagelen int
	cmd := &cobra.Command{
		Use:   "list-comments",
		Short: "List comments on a pull request",
		Long:  "List comments for a pull request and optionally include resolved comments.",
		Example: `bbmd pr list-comments --repo-owner <workspace> --repo-name <repo> --pr-id <id>
bbmd pr list-comments --repo-owner <workspace> --repo-name <repo> --pr-id <id> --include-resolved true`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Flags().Changed("page") && page < 0 {
				return fmt.Errorf("invalid --page: must be non-negative, got %d", page)
			}
			if cmd.Flags().Changed("pagelen") && pagelen < 0 {
				return fmt.Errorf("invalid --pagelen: must be non-negative, got %d", pagelen)
			}
			return runPRListComments(
				cmd,
				container,
				rootParams,
				prListCommentsCLI{
					Params: app.BitbucketListPRCommentsParams{
						RepoOwner:     core.RepoOwner,
						RepoName:      core.RepoName,
						PullRequestID: core.PRID,
						AccountName:   core.Account,
						Page:          page,
						PageLen:       pagelen,
					},
					IncludeResolved: includeResolved,
				},
			)
		},
	}
	bindPRCoreFlags(cmd, &core)
	requirePRCoreFlags(cmd)
	cmd.Flags().BoolVar(&includeResolved, "include-resolved", false, "Include resolved comments in the output")
	cmd.Flags().IntVar(&page, "page", 0, "Page number to retrieve (optional, defaults to first page when omitted)")
	cmd.Flags().IntVar(&pagelen, "pagelen", 0, "Number of comments per page (optional, defaults to 100)")
	return cmd
}

func runPRListComments(
	cmd *cobra.Command,
	container *dig.Container,
	rootParams *rootCommandParams,
	cli prListCommentsCLI,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[prListCommentsCLI, *app.BitbucketListPRCommentsResult]{
			rootParams: rootParams,
			params:     cli,
			target: func(ctx context.Context, p prListCommentsCLI) (*app.BitbucketListPRCommentsResult, error) {
				result, err := svc.ListPRComments(ctx, p.Params)
				if err != nil {
					return nil, err
				}
				applyPRCommentsResolvedFilter(result, p.IncludeResolved)
				return result, nil
			},
		})
	})
}

func newPRResolveCommentCmd(container *dig.Container, rootParams *rootCommandParams) *cobra.Command {
	var core prCoreIDs
	var commentID int64
	cmd := &cobra.Command{
		Use:   "resolve-comment",
		Short: "Resolve a pull request comment thread",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPRResolveComment(cmd, container, rootParams, app.BitbucketResolvePRCommentParams{
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
	rootParams *rootCommandParams,
	params app.BitbucketResolvePRCommentParams,
) error {
	return container.Invoke(func(deps execDeps, svc *app.BitbucketService) error {
		return execAndWrite(cmd, deps, execArgs[app.BitbucketResolvePRCommentParams, *bitbucket.CommentResolution]{
			rootParams: rootParams,
			params:     params,
			target:     svc.ResolvePRComment,
		})
	})
}
