package controllers

import (
	"context"
	"time"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
)

// bitbucketService defines the operations required by the BitbucketController.
// This interface matches the methods from app.BitbucketService that are used by the controller.
type bitbucketService interface {
	CreatePR(ctx context.Context, params app.BitbucketCreatePRParams) (*bitbucket.PullRequest, error)
	ReadPR(ctx context.Context, params app.BitbucketReadPRParams) (*bitbucket.PullRequest, error)
	UpdatePR(ctx context.Context, params app.BitbucketUpdatePRParams) (*bitbucket.PullRequest, error)
	ApprovePR(ctx context.Context, params app.BitbucketApprovePRParams) (*bitbucket.Participant, error)
	MergePR(ctx context.Context, params app.BitbucketMergePRParams) (*bitbucket.PullRequest, error)
	ListTasks(ctx context.Context, params app.BitbucketListTasksParams) (*bitbucket.PaginatedTasks, error)
	UpdateTask(ctx context.Context, params app.BitbucketUpdateTaskParams) (*bitbucket.PullRequestCommentTask, error)
	CreateTask(ctx context.Context, params app.BitbucketCreateTaskParams) (*bitbucket.PullRequestCommentTask, error)
	GetPRDiffStat(ctx context.Context, params app.BitbucketGetPRDiffStatParams) (*app.PaginatedDiffStat, error)
	GetPRDiff(ctx context.Context, params app.BitbucketGetPRDiffParams) (string, error)
	GetFileContent(ctx context.Context, params app.BitbucketGetFileContentParams) (*bitbucket.FileContentResult, error)
	AddPRComment(ctx context.Context, params app.BitbucketAddPRCommentParams) (int64, string, error)
	RequestPRChanges(ctx context.Context, params app.BitbucketRequestPRChangesParams) (string, time.Time, error)
	ListPRComments(
		ctx context.Context,
		params app.BitbucketListPRCommentsParams,
	) (*bitbucket.ListPRCommentsResponse, error)
}

// Ensure that app.BitbucketService implements bitbucketService.
var _ bitbucketService = (*app.BitbucketService)(nil)
