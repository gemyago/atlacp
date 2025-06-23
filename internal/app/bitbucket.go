package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"go.uber.org/dig"
)

// BitbucketService provides business logic for Bitbucket operations.
type BitbucketService struct {
	client       BitbucketClient
	accountsRepo AtlassianAccountsRepository
	logger       *slog.Logger
}

// BitbucketServiceDeps contains dependencies for the Bitbucket service.
type BitbucketServiceDeps struct {
	dig.In

	Client       BitbucketClient
	AccountsRepo AtlassianAccountsRepository
	RootLogger   *slog.Logger
}

// NewBitbucketService creates a new Bitbucket service.
func NewBitbucketService(deps BitbucketServiceDeps) *BitbucketService {
	return &BitbucketService{
		client:       deps.Client,
		accountsRepo: deps.AccountsRepo,
		logger:       deps.RootLogger.WithGroup("app.bitbucket-service"),
	}
}

// BitbucketCreatePRParams contains parameters for creating a pull request.
type BitbucketCreatePRParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Title of the pull request
	Title string `json:"title"`

	// Description of the pull request
	Description string `json:"description"`

	// Source branch name
	SourceBranch string `json:"source_branch"`

	// Destination branch name
	DestBranch string `json:"dest_branch"`

	// Whether to close the source branch after merging
	CloseSourceBranch bool `json:"close_source_branch"`

	// Reviewer usernames (optional)
	Reviewers []string `json:"reviewers,omitempty"`
}

// BitbucketReadPRParams contains parameters for retrieving a pull request.
type BitbucketReadPRParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`
}

// BitbucketUpdatePRParams contains parameters for updating a pull request.
type BitbucketUpdatePRParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`

	// Updated title (optional)
	Title string `json:"title,omitempty"`

	// Updated description (optional)
	Description string `json:"description,omitempty"`
}

// BitbucketApprovePRParams contains parameters for approving a pull request.
type BitbucketApprovePRParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`
}

// BitbucketMergePRParams contains parameters for merging a pull request.
type BitbucketMergePRParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`

	// Merge commit message (optional)
	Message string `json:"message,omitempty"`

	// Whether to close the source branch after merging
	CloseSourceBranch bool `json:"close_source_branch"`

	// Merge strategy (merge_commit, squash, fast_forward)
	MergeStrategy string `json:"merge_strategy,omitempty"`
}

// CreatePR creates a new pull request.
func (s *BitbucketService) CreatePR(
	ctx context.Context,
	params BitbucketCreatePRParams,
) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Creating pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.String("source", params.SourceBranch),
		slog.String("dest", params.DestBranch))

	return nil, errors.New("not implemented")
}

// ReadPR retrieves a specific pull request.
func (s *BitbucketService) ReadPR(ctx context.Context, params BitbucketReadPRParams) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Reading pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	return nil, errors.New("not implemented")
}

// UpdatePR updates an existing pull request.
func (s *BitbucketService) UpdatePR(
	ctx context.Context,
	params BitbucketUpdatePRParams,
) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Updating pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	return nil, errors.New("not implemented")
}

// ApprovePR approves a pull request.
func (s *BitbucketService) ApprovePR(
	ctx context.Context,
	params BitbucketApprovePRParams,
) (*bitbucket.Participant, error) {
	s.logger.InfoContext(ctx, "Approving pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	return nil, errors.New("not implemented")
}

// MergePR merges a pull request.
func (s *BitbucketService) MergePR(ctx context.Context, params BitbucketMergePRParams) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Merging pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID),
		slog.String("strategy", params.MergeStrategy))

	return nil, errors.New("not implemented")
}
