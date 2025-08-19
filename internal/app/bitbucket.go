package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"go.uber.org/dig"
)

// Task states.
const (
	// TaskStateResolved is the state value for a resolved task.
	TaskStateResolved = "RESOLVED"

	// TaskStateUnresolved is the state value for an unresolved task.
	TaskStateUnresolved = "UNRESOLVED"
)

// BitbucketService provides business logic for Bitbucket operations.
type BitbucketService struct {
	client      bitbucketClient
	authFactory bitbucketAuthFactory
	logger      *slog.Logger
}

// BitbucketServiceDeps contains dependencies for the Bitbucket service.
type BitbucketServiceDeps struct {
	dig.In

	Client      bitbucketClient
	AuthFactory bitbucketAuthFactory
	RootLogger  *slog.Logger
}

// NewBitbucketService creates a new Bitbucket service.
func NewBitbucketService(deps BitbucketServiceDeps) *BitbucketService {
	return &BitbucketService{
		client:      deps.Client,
		authFactory: deps.AuthFactory,
		logger:      deps.RootLogger.WithGroup("app.bitbucket-service"),
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

	// Whether to create the pull request as a draft
	Draft *bool `json:"draft,omitempty"`
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

	// Whether to update the pull request as a draft
	Draft *bool `json:"draft,omitempty"`
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

type BitbucketRequestPRChangesParams struct {
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

// BitbucketListTasksParams contains parameters for listing tasks on a pull request.
type BitbucketListTasksParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`

	// Optional query to filter tasks (optional)
	Query string `json:"query,omitempty"`

	// Sort order for tasks (optional)
	Sort string `json:"sort,omitempty"`

	// Maximum number of tasks to return per page (optional)
	PageLen int `json:"page_len,omitempty"`
}

// BitbucketUpdateTaskParams contains parameters for updating a task on a pull request.
type BitbucketUpdateTaskParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`

	// Task ID to update
	TaskID int `json:"task_id"`

	// Updated task content (optional)
	Content string `json:"content,omitempty"`

	// Updated task state: "RESOLVED" or "UNRESOLVED" (optional)
	State string `json:"state,omitempty"`
}

// BitbucketCreateTaskParams contains parameters for creating a task on a pull request.
type BitbucketCreateTaskParams struct {
	// Account name to use for authentication (optional, uses default if empty)
	AccountName string `json:"account_name,omitempty"`

	// Repository owner (username/workspace)
	RepoOwner string `json:"repo_owner"`

	// Repository name (slug)
	RepoName string `json:"repo_name"`

	// Pull request ID
	PullRequestID int `json:"pull_request_id"`

	// Task content
	Content string `json:"content"`

	// Comment ID to associate with the task (optional)
	CommentID int64 `json:"comment_id,omitempty"`

	// Task state: "RESOLVED" or "UNRESOLVED" (optional)
	// If not provided, defaults to "UNRESOLVED"
	State string `json:"state,omitempty"`
}

// CreatePR creates a new pull request.
func (s *BitbucketService) CreatePR(
	ctx context.Context,
	params BitbucketCreatePRParams,
) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Creating pull request",
		slog.String("repo", params.RepoName),
		slog.String("source", params.SourceBranch),
		slog.String("dest", params.DestBranch))

	// Validate required parameters
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.Title == "" {
		return nil, errors.New("title is required")
	}
	if params.SourceBranch == "" {
		return nil, errors.New("source branch is required")
	}
	if params.DestBranch == "" {
		return nil, errors.New("destination branch is required")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Build pull request object
	prRequest := &bitbucket.PullRequest{
		Title:             params.Title,
		Description:       params.Description,
		CloseSourceBranch: params.CloseSourceBranch,
		Source: bitbucket.PullRequestSource{
			Branch: bitbucket.PullRequestBranch{
				Name: params.SourceBranch,
			},
		},
		Destination: &bitbucket.PullRequestDestination{
			Branch: bitbucket.PullRequestBranch{
				Name: params.DestBranch,
			},
		},
		Draft: params.Draft,
	}

	// Add reviewers if specified
	if len(params.Reviewers) > 0 {
		prRequest.Reviewers = make([]bitbucket.PullRequestAuthor, len(params.Reviewers))
		for i, reviewer := range params.Reviewers {
			prRequest.Reviewers[i] = bitbucket.PullRequestAuthor{
				Username: reviewer,
			}
		}
	}

	// Call the client to create the pull request
	return s.client.CreatePR(ctx, tokenProvider, bitbucket.CreatePRParams{
		Username: params.RepoOwner,
		RepoSlug: params.RepoName,
		Request:  prRequest,
	})
}

// ReadPR retrieves a specific pull request.
func (s *BitbucketService) ReadPR(ctx context.Context, params BitbucketReadPRParams) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Reading pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Call the client to get the pull request
	pr, err := s.client.GetPR(ctx, tokenProvider, bitbucket.GetPRParams{
		Username:      params.RepoOwner,
		RepoSlug:      params.RepoName,
		PullRequestID: params.PullRequestID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return pr, nil
}

// UpdatePR updates an existing pull request.
func (s *BitbucketService) UpdatePR(
	ctx context.Context,
	params BitbucketUpdatePRParams,
) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Updating pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}
	if params.Title == "" && params.Description == "" && params.Draft == nil {
		return nil, errors.New("either title, description or draft must be provided")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Create update request with provided fields
	updateRequest := &bitbucket.PullRequest{
		Title:       params.Title,
		Description: params.Description,
		Draft:       params.Draft,
	}

	// Call the client to update the pull request
	pr, err := s.client.UpdatePR(ctx, tokenProvider, bitbucket.UpdatePRParams{
		Username:      params.RepoOwner,
		RepoSlug:      params.RepoName,
		PullRequestID: params.PullRequestID,
		Request:       updateRequest,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update pull request: %w", err)
	}

	return pr, nil
}

// ApprovePR approves a pull request.
func (s *BitbucketService) ApprovePR(
	ctx context.Context,
	params BitbucketApprovePRParams,
) (*bitbucket.Participant, error) {
	s.logger.InfoContext(ctx, "Approving pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Call the client to approve the pull request
	participant, err := s.client.ApprovePR(ctx, tokenProvider, bitbucket.ApprovePRParams{
		Username:      params.RepoOwner,
		RepoSlug:      params.RepoName,
		PullRequestID: params.PullRequestID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to approve pull request: %w", err)
	}

	return participant, nil
}

// RequestPRChanges requests changes on a pull request.
func (s *BitbucketService) RequestPRChanges(
	ctx context.Context,
	params BitbucketRequestPRChangesParams,
) (string, time.Time, error) {
	s.logger.InfoContext(ctx, "Requesting changes on pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID),
	)

	// Validate required parameters
	if params.RepoOwner == "" {
		return "", time.Time{}, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return "", time.Time{}, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return "", time.Time{}, errors.New("pull request ID must be positive")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Call the client to request PR changes
	status, ts, err := s.client.RequestPRChanges(ctx, tokenProvider, bitbucket.RequestPRChangesParams{
		Workspace: params.RepoOwner,
		RepoSlug:  params.RepoName,
		PullReqID: params.PullRequestID,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return status, ts, nil
}

// MergePR merges a pull request.
func (s *BitbucketService) MergePR(ctx context.Context, params BitbucketMergePRParams) (*bitbucket.PullRequest, error) {
	s.logger.InfoContext(ctx, "Merging pull request",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID),
		slog.String("strategy", params.MergeStrategy))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}

	// Validate merge strategy if provided
	if params.MergeStrategy != "" && !isValidMergeStrategy(params.MergeStrategy) {
		return nil, errors.New("invalid merge strategy: must be one of merge_commit, squash, or fast_forward")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Create merge parameters
	mergeParams := &bitbucket.PullRequestMergeParameters{
		CloseSourceBranch: params.CloseSourceBranch,
		Message:           params.Message,
	}

	// Only add merge strategy if specified
	if params.MergeStrategy != "" {
		mergeParams.MergeStrategy = params.MergeStrategy
	}

	// Call the client to merge the pull request
	pr, err := s.client.MergePR(ctx, tokenProvider, bitbucket.MergePRParams{
		Username:        params.RepoOwner,
		RepoSlug:        params.RepoName,
		PullRequestID:   params.PullRequestID,
		MergeParameters: mergeParams,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to merge pull request: %w", err)
	}

	return pr, nil
}

// isValidMergeStrategy checks if the provided merge strategy is valid.
func isValidMergeStrategy(strategy string) bool {
	validStrategies := map[string]bool{
		"merge_commit": true,
		"squash":       true,
		"fast_forward": true,
		"":             true, // Empty is valid, will use repo default
	}
	return validStrategies[strategy]
}

// ListTasks retrieves a list of tasks for a specific pull request.
func (s *BitbucketService) ListTasks(
	ctx context.Context,
	params BitbucketListTasksParams,
) (*bitbucket.PaginatedTasks, error) {
	s.logger.InfoContext(ctx, "Listing pull request tasks",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Call the client to list tasks
	tasks, err := s.client.ListPullRequestTasks(ctx, tokenProvider, bitbucket.ListPullRequestTasksParams{
		Workspace: params.RepoOwner,
		RepoSlug:  params.RepoName,
		PullReqID: params.PullRequestID,
		Query:     params.Query,
		Sort:      params.Sort,
		PageLen:   params.PageLen,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pull request tasks: %w", err)
	}

	return tasks, nil
}

// UpdateTask updates a task on a pull request.
func (s *BitbucketService) UpdateTask(
	ctx context.Context,
	params BitbucketUpdateTaskParams,
) (*bitbucket.PullRequestCommentTask, error) {
	s.logger.InfoContext(ctx, "Updating pull request task",
		slog.String("repo", params.RepoOwner+"/"+params.RepoName),
		slog.Int("pr_id", params.PullRequestID),
		slog.Int("task_id", params.TaskID))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}
	if params.TaskID <= 0 {
		return nil, errors.New("task ID must be positive")
	}
	if params.Content == "" && params.State == "" {
		return nil, errors.New("either content or state must be provided")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Call the client to update the task
	task, err := s.client.UpdateTask(ctx, tokenProvider, bitbucket.UpdateTaskParams{
		Workspace: params.RepoOwner,
		RepoSlug:  params.RepoName,
		PullReqID: params.PullRequestID,
		TaskID:    params.TaskID,
		Content:   params.Content,
		State:     params.State,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// CreateTask creates a new task on a pull request.
func (s *BitbucketService) CreateTask(
	ctx context.Context,
	params BitbucketCreateTaskParams,
) (*bitbucket.PullRequestCommentTask, error) {
	s.logger.InfoContext(ctx, "Creating task on pull request",
		slog.String("repo", params.RepoName),
		slog.Int("pr_id", params.PullRequestID))

	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}
	if params.Content == "" {
		return nil, errors.New("content is required")
	}

	// Get token provider from auth factory
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

	// Prepare the client parameters
	clientParams := bitbucket.CreatePullRequestTaskParams{
		Workspace: params.RepoOwner,
		RepoSlug:  params.RepoName,
		PullReqID: params.PullRequestID,
		Content:   params.Content,
		CommentID: params.CommentID,
	}

	// Handle optional state parameter
	if params.State != "" {
		// Convert state to pending flag (RESOLVED -> false, UNRESOLVED -> true)
		pending := params.State != TaskStateResolved
		clientParams.Pending = &pending
	}

	// Call the client to create the task
	task, err := s.client.CreatePullRequestTask(ctx, tokenProvider, clientParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

type PaginatedDiffStat struct {
	Size    int                  `json:"size,omitempty"`
	Page    int                  `json:"page,omitempty"`
	PageLen int                  `json:"pagelen,omitempty"`
	Values  []bitbucket.DiffStat `json:"values"`
}

type BitbucketGetPRDiffStatParams struct {
	AccountName   string
	RepoOwner     string
	RepoName      string
	PullRequestID int
}

func (s *BitbucketService) GetPRDiffStat(
	ctx context.Context,
	params BitbucketGetPRDiffStatParams,
) (*PaginatedDiffStat, error) {
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)
	clientResult, err := s.client.GetPRDiffStat(
		ctx,
		tokenProvider,
		bitbucket.GetPRDiffStatParams{
			RepoOwner: params.RepoOwner,
			RepoName:  params.RepoName,
			PRID:      params.PullRequestID,
		},
	)
	if err != nil {
		return nil, err
	}
	return &PaginatedDiffStat{
		Size:    clientResult.Size,
		Page:    clientResult.Page,
		PageLen: clientResult.PageLen,
		Values:  clientResult.Values,
	}, nil
}

type BitbucketGetPRDiffParams struct {
	AccountName   string
	RepoOwner     string
	RepoName      string
	PullRequestID int
	FilePaths     []string
	ContextLines  *int
}

func (s *BitbucketService) GetPRDiff(
	ctx context.Context,
	params BitbucketGetPRDiffParams,
) (*bitbucket.Diff, error) {
	// Validate required parameters
	if params.RepoOwner == "" {
		return nil, errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return nil, errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return nil, errors.New("pull request ID must be positive")
	}

	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)
	clientParams := bitbucket.GetPRDiffParams{
		RepoOwner: params.RepoOwner,
		RepoName:  params.RepoName,
		PRID:      params.PullRequestID,
		FilePaths: params.FilePaths,
		Context:   params.ContextLines,
	}
	return s.client.GetPRDiff(ctx, tokenProvider, clientParams)
}

type BitbucketGetFileContentParams struct {
	AccountName string
	RepoOwner   string
	RepoName    string
	Commit      string
	Path        string
}

func (s *BitbucketService) GetFileContent(
	ctx context.Context,
	params BitbucketGetFileContentParams,
) (*bitbucket.FileContentResult, error) {
	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)
	clientParams := bitbucket.GetFileContentParams{
		RepoOwner:  params.RepoOwner,
		RepoName:   params.RepoName,
		CommitHash: params.Commit,
		FilePath:   params.Path,
	}
	fileContent, err := s.client.GetFileContent(ctx, tokenProvider, clientParams)
	if err != nil {
		return nil, err
	}
	// Minimal stub: meta fields are hardcoded for now
	return &bitbucket.FileContentResult{
		Content: fileContent.Content,
		Meta: bitbucket.FileContentMeta{
			Size:     len(fileContent.Content),
			Type:     "file",
			Encoding: "utf-8",
		},
	}, nil
}

type BitbucketAddPRCommentParams struct {
	AccountName   string
	RepoOwner     string
	RepoName      string
	PullRequestID int
	Content       string
	FilePath      string
	LineFrom      int
	LineTo        int
}

// AddPRComment adds a comment to a pull request (general or inline).
func (s *BitbucketService) AddPRComment(
	ctx context.Context,
	params BitbucketAddPRCommentParams,
) (int64, string, error) {
	// Validate required parameters
	if params.RepoOwner == "" {
		return 0, "", errors.New("repository owner is required")
	}
	if params.RepoName == "" {
		return 0, "", errors.New("repository name is required")
	}
	if params.PullRequestID <= 0 {
		return 0, "", errors.New("pull request ID must be positive")
	}
	if params.Content == "" {
		return 0, "", errors.New("comment content is required")
	}

	tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)
	clientParams := bitbucket.AddPRCommentParams{
		Workspace:   params.RepoOwner,
		RepoSlug:    params.RepoName,
		PullReqID:   params.PullRequestID,
		CommentText: params.Content,
		FilePath:    params.FilePath,
		LineFrom:    params.LineFrom,
		LineTo:      params.LineTo,
		Account:     params.AccountName,
	}
	return s.client.AddPRComment(ctx, tokenProvider, clientParams)
}
