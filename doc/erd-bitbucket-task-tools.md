# Engineering Requirements Document: Bitbucket PR Task Tools

## Introduction/Overview
This document outlines the engineering requirements for implementing two new Bitbucket MCP tools to manage pull request tasks:
1. List pull request tasks
2. Update pull request task
3. Create pull request task

These tools will extend the existing Bitbucket MCP controller to provide task management capabilities, allowing users to view, update and create tasks associated with pull requests.

## Business Logic
The implementation will follow the existing pattern in the codebase where:
1. The MCP controller exposes tools that accept parameters from the user
2. The controller validates inputs and transforms them into service layer parameters
3. The controller calls the appropriate service methods in the app layer
4. The app layer interacts with the client layer to perform the operations
5. The controller formats the response for the user

## High Level Architecture
The implementation will require changes at multiple layers:

1. **App Layer**:
   - Add new parameter structs for task operations
   - Add new methods to the BitbucketService
   - Update the bitbucketClient interface to include task operations

2. **Controller Layer**:
   - Add two new tools to the BitbucketController:
     - `bitbucket_list_pr_tasks` - Lists tasks on a pull request
     - `bitbucket_update_pr_task` - Updates a task on a pull request

These tools will leverage the existing client methods through the app layer:
- `ListPullRequestTasks` in `internal/services/bitbucket/list_tasks.go`
- `UpdateTask` in `internal/services/bitbucket/update_task.go`

## Detailed Architecture

### App Layer Changes

#### 1. Update bitbucketClient interface in ports.go
Add the following methods to the interface:
```go
// ListPullRequestTasks lists tasks on a pull request.
ListPullRequestTasks(
    ctx context.Context,
    tokenProvider bitbucket.TokenProvider,
    params bitbucket.ListPullRequestTasksParams,
) (*bitbucket.PaginatedTasks, error)

// UpdateTask updates a task on a pull request.
UpdateTask(
    ctx context.Context,
    tokenProvider bitbucket.TokenProvider,
    params bitbucket.UpdateTaskParams,
) (*bitbucket.PullRequestCommentTask, error)

// CreateTask creates a new task on a pull request.
CreateTask(
    ctx context.Context,
    tokenProvider bitbucket.TokenProvider,
    params bitbucket.CreateTaskParams,
) (*bitbucket.PullRequestCommentTask, error)
```

#### 2. Add new parameter structs in bitbucket.go
```go
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

    // Query filter for tasks (optional)
    Query string `json:"query,omitempty"`

    // Sort order for tasks (optional)
    Sort string `json:"sort,omitempty"`

    // Number of tasks per page (optional)
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

    // Task ID
    TaskID int `json:"task_id"`

    // Updated task content (optional)
    Content string `json:"content,omitempty"`

    // Task state (optional, "RESOLVED" or "UNRESOLVED")
    State string `json:"state,omitempty"`
}

// BitbucketCreateTaskParams contains parameters for creating a new task on a pull request.
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
    
    // Comment ID to add task to (optional, if not provided task will be created as a standalone task)
    CommentID int `json:"comment_id,omitempty"`
}
```

#### 3. Add new methods to BitbucketService
```go
// ListTasks lists tasks on a pull request.
func (s *BitbucketService) ListTasks(
    ctx context.Context,
    params BitbucketListTasksParams,
) (*bitbucket.PaginatedTasks, error) {
    s.logger.InfoContext(ctx, "Listing tasks on pull request",
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
        return nil, fmt.Errorf("failed to list tasks: %w", err)
    }

    return tasks, nil
}

// UpdateTask updates a task on a pull request.
func (s *BitbucketService) UpdateTask(
    ctx context.Context,
    params BitbucketUpdateTaskParams,
) (*bitbucket.PullRequestCommentTask, error) {
    s.logger.InfoContext(ctx, "Updating task on pull request",
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
    if params.State != "" && params.State != "RESOLVED" && params.State != "UNRESOLVED" {
        return nil, errors.New("state must be either RESOLVED or UNRESOLVED")
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
    if params.Content == "" {
        return nil, errors.New("task content is required")
    }

    // Get token provider from auth factory
    tokenProvider := s.authFactory.getTokenProvider(ctx, params.AccountName)

    // Call the client to create the task
    task, err := s.client.CreateTask(ctx, tokenProvider, bitbucket.CreateTaskParams{
        Workspace: params.RepoOwner,
        RepoSlug:  params.RepoName,
        PullReqID: params.PullRequestID,
        Content:   params.Content,
        CommentID: params.CommentID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create task: %w", err)
    }

    return task, nil
}
```

### Controller Layer Changes

#### 1. Update bitbucketService interface in the controller
Add the new methods to the interface:
```go
ListTasks(ctx context.Context, params app.BitbucketListTasksParams) (*bitbucket.PaginatedTasks, error)
UpdateTask(ctx context.Context, params app.BitbucketUpdateTaskParams) (*bitbucket.PullRequestCommentTask, error)
CreateTask(ctx context.Context, params app.BitbucketCreateTaskParams) (*bitbucket.PullRequestCommentTask, error)
```

#### 2. List PR Tasks Tool
1. **Tool Name**: `bitbucket_list_pr_tasks`
2. **Parameters**:
   - Required:
     - `repo_owner` (string): Repository owner/workspace
     - `repo_name` (string): Repository name (slug)
     - `pr_id` (number): Pull request ID
   - Optional:
     - `account` (string): Atlassian account name to use
     - `query` (string): Query filter for tasks
     - `sort` (string): Sort order for tasks
     - `page_len` (number): Number of tasks per page
3. **Response**:
   - Summary text with task count
   - Full task data as JSON

#### 3. Update PR Task Tool
1. **Tool Name**: `bitbucket_update_pr_task`
2. **Parameters**:
   - Required:
     - `repo_owner` (string): Repository owner/workspace
     - `repo_name` (string): Repository name (slug)
     - `pr_id` (number): Pull request ID
     - `task_id` (number): Task ID
   - Optional:
     - `account` (string): Atlassian account name to use
     - `content` (string): Updated task content
     - `state` (string): Task state ("RESOLVED" or "UNRESOLVED")
3. **Response**:
   - Confirmation text with task update details

#### 4. Create PR Task Tool
1. **Tool Name**: `bitbucket_create_pr_task`
2. **Parameters**:
   - Required:
     - `repo_owner` (string): Repository owner/workspace
     - `repo_name` (string): Repository name (slug)
     - `pr_id` (number): Pull request ID
     - `content` (string): Task content
   - Optional:
     - `account` (string): Atlassian account name to use
     - `comment_id` (number): Comment ID to add task to
3. **Response**:
   - Confirmation text with task creation details
   - Task ID and state in the response

### Implementation Details
1. Add new methods to the app layer:
   - `ListTasks` in `BitbucketService`
   - `UpdateTask` in `BitbucketService`
   - `CreateTask` in `BitbucketService`
2. Update the `bitbucketClient` interface to include task-related methods
3. Add new controller methods:
   - `newListPRTasksServerTool()`
   - `newUpdatePRTaskServerTool()`
   - `newCreatePRTaskServerTool()`
4. Add these tools to the `NewTools()` method
5. Follow the same pattern as existing tools for parameter validation and error handling

## Key Architectural Decisions
1. **Consistent Layering**: Follow the existing pattern of controller → app service → client
2. **Parameter Validation**: Validate parameters at both the controller and app layers
3. **Response Format**: The list tasks tool will return both summary text and full JSON data, similar to the read PR tool
4. **Error Handling**: Use the same error handling pattern as existing tools, returning `mcp.NewToolResultError` for validation errors
5. **Interface Extension**: Extend the `bitbucketClient` interface to include the task-related methods

## Testing Strategy
1. **Unit Tests**: Add unit tests for the new app layer methods
   - Test parameter validation
   - Test successful execution with mocked client
   - Test error handling
2. **Integration Testing**: Test the MCP tools manually to verify end-to-end functionality
3. **Test Coverage**: Ensure all code paths are covered, including error cases

## Open Questions
1. Should we add more validation for optional parameters like `state`? - no
2. Should we support pagination for task listing beyond the initial page? - no
3. Should we add methods for creating and deleting tasks in the future? - maybe
4. Should we support adding tasks to specific comment threads besides creating standalone tasks? - yes