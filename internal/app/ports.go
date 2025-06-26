package app

import (
	"context"
	"errors"

	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/gemyago/atlacp/internal/services/http/middleware"
)

// AtlassianAccountsRepository defines the port for accessing Atlassian account information.
// This is an outbound port that will be implemented by the infrastructure layer.
type AtlassianAccountsRepository interface {
	// GetDefaultAccount returns the default Atlassian account configuration.
	// Returns an error if no default account is found.
	GetDefaultAccount(ctx context.Context) (*AtlassianAccount, error)

	// GetAccountByName returns an account with the specified name.
	// Returns an error if no account with the name is found.
	GetAccountByName(ctx context.Context, name string) (*AtlassianAccount, error)
}

// TokenProvider provides authentication tokens for API requests.
type TokenProvider interface {
	// GetToken returns an authentication token for API requests.
	GetToken(ctx context.Context) (middleware.Token, error)
}

// bitbucketClient defines the interface for Bitbucket API operations.
// This is an outbound port that will be implemented by the infrastructure layer.
type bitbucketClient interface {
	// CreatePR creates a new pull request in the specified repository.
	CreatePR(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.CreatePRParams,
	) (*bitbucket.PullRequest, error)

	// GetPR retrieves a specific pull request by ID.
	GetPR(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.GetPRParams,
	) (*bitbucket.PullRequest, error)

	// UpdatePR updates an existing pull request.
	UpdatePR(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.UpdatePRParams,
	) (*bitbucket.PullRequest, error)

	// ApprovePR approves a pull request.
	ApprovePR(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.ApprovePRParams,
	) (*bitbucket.Participant, error)

	// MergePR merges a pull request.
	MergePR(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.MergePRParams,
	) (*bitbucket.PullRequest, error)

	// ListPullRequestTasks returns a paginated list of tasks on a pull request.
	ListPullRequestTasks(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.ListPullRequestTasksParams,
	) (*bitbucket.PaginatedTasks, error)

	// UpdateTask updates an existing task on a pull request.
	UpdateTask(
		ctx context.Context,
		tokenProvider bitbucket.TokenProvider,
		params bitbucket.UpdateTaskParams,
	) (*bitbucket.PullRequestCommentTask, error)
}

// Error types for account-related operations.
var (
	// ErrNoDefaultAccount is returned when no default account is configured.
	ErrNoDefaultAccount = errors.New("no default Atlassian account configured")

	// ErrAccountNotFound is returned when a specific named account is not found.
	ErrAccountNotFound = errors.New("atlassian account not found")

	// ErrAccountConfigInvalid is returned when account configuration is invalid.
	ErrAccountConfigInvalid = errors.New("atlassian account configuration is invalid")
)
