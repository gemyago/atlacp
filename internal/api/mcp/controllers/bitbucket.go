package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/samber/lo"
	"go.uber.org/dig"
)

const (
	// TaskStateResolved represents a "RESOLVED" task state.
	TaskStateResolved = "RESOLVED"

	// TaskStateUnresolved represents an "UNRESOLVED" task state.
	TaskStateUnresolved = "UNRESOLVED"
)

// BitbucketControllerDeps contains dependencies for the Bitbucket MCP controller.
type BitbucketControllerDeps struct {
	dig.In

	RootLogger       *slog.Logger
	BitbucketService bitbucketService
}

// BitbucketController provides MCP Bitbucket tool functionality.
type BitbucketController struct {
	logger           *slog.Logger
	bitbucketService bitbucketService
}

// NewBitbucketController creates a new Bitbucket MCP controller.
func NewBitbucketController(deps BitbucketControllerDeps) *BitbucketController {
	return &BitbucketController{
		logger:           deps.RootLogger.WithGroup("mcp.bitbucket-controller"),
		bitbucketService: deps.BitbucketService,
	}
}

// newCreatePRServerTool returns a server tool for creating pull requests.
func (bc *BitbucketController) newCreatePRServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_create_pr",
		mcp.WithDescription("Create a pull request in Bitbucket"),
		mcp.WithString("title",
			mcp.Description("Pull request title"),
			mcp.Required(),
		),
		mcp.WithString("source_branch",
			mcp.Description("Source branch name"),
			mcp.Required(),
		),
		mcp.WithString("target_branch",
			mcp.Description("Target branch name"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("description",
			mcp.Description("Pull request description"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
		mcp.WithBoolean("draft",
			mcp.Description("Create as draft pull request (optional, defaults to false)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_create_pr request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid title parameter", err), nil
		}

		sourceBranch, err := request.RequireString("source_branch")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid source_branch parameter", err), nil
		}

		targetBranch, err := request.RequireString("target_branch")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid target_branch parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		description := request.GetString("description", "")
		account := request.GetString("account", "")
		draft := request.GetBool("draft", false)

		// Create parameters for the service layer
		params := app.BitbucketCreatePRParams{
			Title:        title,
			SourceBranch: sourceBranch,
			DestBranch:   targetBranch,
			Description:  description,
			AccountName:  account,
			RepoOwner:    repoOwner,
			RepoName:     repoName,
			Draft:        lo.ToPtr(draft),
		}

		// Call the service to create the pull request
		pr, err := bc.bitbucketService.CreatePR(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create pull request: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Created pull request #%d: %s", pr.ID, pr.Title)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newReadPRServerTool returns a server tool for reading pull request details.
func (bc *BitbucketController) newReadPRServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_read_pr",
		mcp.WithDescription("Get pull request details from Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_read_pr request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		account := request.GetString("account", "")

		// Create parameters for the service layer
		params := app.BitbucketReadPRParams{
			PullRequestID: prID,
			AccountName:   account,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
		}

		// Call the service to read the pull request
		pr, err := bc.bitbucketService.ReadPR(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to read pull request: %w", err)
		}

		// Convert PR to JSON for the resource
		prJSON, err := json.MarshalIndent(pr, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal pull request to JSON: %w", err)
		}

		// Create a summary text for the PR
		summaryText := fmt.Sprintf("Pull request #%d: %s (Status: %s)", pr.ID, pr.Title, pr.State)

		// Return both a summary text and the full PR data as a resource
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: summaryText,
				},

				// Sending json as text since some clients (Cursor)
				// do not support resources (at least not yet)
				mcp.NewTextContent(string(prJSON)),
			},
		}, nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newUpdatePRServerTool returns a server tool for updating pull requests.
func (bc *BitbucketController) newUpdatePRServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_update_pr",
		mcp.WithDescription("Update a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("title",
			mcp.Description("New pull request title"),
		),
		mcp.WithString("description",
			mcp.Description("New pull request description"),
		),
		mcp.WithBoolean("draft",
			mcp.Description("Update as draft pull request (optional)"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_update_pr request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// At least one of title, description, or draft must be provided
		title := request.GetString("title", "")
		description := request.GetString("description", "")

		allArgs := request.GetArguments()
		_, hasTitle := allArgs["title"]
		_, hasDescription := allArgs["description"]
		_, hasDraft := allArgs["draft"]

		if !hasTitle && !hasDescription && !hasDraft {
			return mcp.NewToolResultError("Missing attributes to update a PR"), nil
		}

		var draft *bool
		if hasDraft {
			draft = lo.ToPtr(request.GetBool("draft", false))
		}

		// Optional parameters
		account := request.GetString("account", "")

		// Create parameters for the service layer
		params := app.BitbucketUpdatePRParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			Title:         title,
			Description:   description,
			Draft:         draft,
			AccountName:   account,
		}

		// Call the service to update the pull request
		pr, err := bc.bitbucketService.UpdatePR(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update pull request: %w", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Updated pull request #%d: %s", pr.ID, pr.Title)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newApprovePRServerTool returns a server tool for approving pull requests.
func (bc *BitbucketController) newApprovePRServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_approve_pr",
		mcp.WithDescription("Approve a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_approve_pr request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		account := request.GetString("account", "")

		// Create parameters for the service layer
		params := app.BitbucketApprovePRParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
		}

		// Call the service to approve the pull request
		participant, err := bc.bitbucketService.ApprovePR(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to approve pull request: %w", err)
		}

		// Create a response with the approval details
		return mcp.NewToolResultText(fmt.Sprintf("Pull request #%d approved by %s",
			prID, participant.User.DisplayName)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newMergePRServerTool returns a server tool for merging pull requests.
func (bc *BitbucketController) newMergePRServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_merge_pr",
		mcp.WithDescription("Merge a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("merge_strategy",
			mcp.Description("Merge strategy (merge_commit, squash, fast_forward)"),
		),
		mcp.WithString("commit_message",
			mcp.Description("Custom commit message for the merge (optional)"),
		),
		mcp.WithString("close_source_branch",
			mcp.Description("Whether to close the source branch after merge (true/false)"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_merge_pr request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		mergeStrategy := request.GetString("merge_strategy", "")
		commitMessage := request.GetString("commit_message", "")
		closeSourceBranch := request.GetBool("close_source_branch", false)
		account := request.GetString("account", "")

		// Create parameters for the service layer
		params := app.BitbucketMergePRParams{
			PullRequestID:     prID,
			RepoOwner:         repoOwner,
			RepoName:          repoName,
			MergeStrategy:     mergeStrategy,
			Message:           commitMessage,
			CloseSourceBranch: closeSourceBranch,
			AccountName:       account,
		}

		// Call the service to merge the pull request
		pr, err := bc.bitbucketService.MergePR(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to merge pull request: %w", err)
		}

		// Create a response with the merge details
		var strategyText string
		if mergeStrategy != "" {
			strategyText = fmt.Sprintf(" using %s strategy", mergeStrategy)
		}

		var closeBranchText string
		if closeSourceBranch {
			closeBranchText = " and source branch was closed"
		}

		return mcp.NewToolResultText(fmt.Sprintf("Pull request #%d successfully merged%s%s",
			pr.ID, strategyText, closeBranchText)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newListPRTasksServerTool returns a server tool for listing tasks on a pull request.
func (bc *BitbucketController) newListPRTasksServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_list_pr_tasks",
		mcp.WithDescription("List tasks on a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_list_pr_tasks request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		account := request.GetString("account", "")

		// Create parameters for the service layer
		params := app.BitbucketListTasksParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
		}

		// Call the service to list tasks
		tasks, err := bc.bitbucketService.ListTasks(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks: %w", err)
		}

		// Create a response with the tasks details
		var responseText string
		if tasks.Size == 0 {
			responseText = "No tasks found for this pull request"
		} else {
			responseText = fmt.Sprintf("Found %d tasks", tasks.Size)
			for _, task := range tasks.Values {
				// Get creator display name, handling nil Creator
				var creatorName string
				if task.Creator != nil && task.Creator.DisplayName != "" {
					creatorName = task.Creator.DisplayName
				} else {
					creatorName = "unknown user"
				}

				responseText += fmt.Sprintf("\nTask #%d: [%s] %s (by %s)",
					task.ID,
					task.State,
					task.Content.Raw,
					creatorName)
			}
		}

		return mcp.NewToolResultText(responseText), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newUpdatePRTaskServerTool returns a server tool for updating a task on a pull request.
func (bc *BitbucketController) newUpdatePRTaskServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_update_pr_task",
		mcp.WithDescription("Update a task on a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithNumber("task_id",
			mcp.Description("Task ID to update"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("New content for the task (optional if state is provided)"),
		),
		mcp.WithString("state",
			mcp.Description("New state for the task: RESOLVED or UNRESOLVED (optional if content is provided)"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := bc.makeUpdatePRTaskHandler()

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// makeUpdatePRTaskHandler creates a handler function for the update PR task tool.
// This is split out to reduce the overall function length.
func (bc *BitbucketController) makeUpdatePRTaskHandler() func(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_update_pr_task request", "params", request.Params)

		// Extract parameters and validate
		params, errResult := bc.validateUpdateTaskParams(request)
		if errResult != nil {
			return errResult, nil // This returns a tool result error
		}

		// Call the service to update the task
		task, err := bc.bitbucketService.UpdateTask(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update task: %w", err)
		}

		// Build response text
		responseText := bc.formatUpdateTaskResponse(task, params.TaskID, params.Content, params.State)
		return mcp.NewToolResultText(responseText), nil
	}
}

// validateUpdateTaskParams extracts and validates parameters for the update task request.
func (bc *BitbucketController) validateUpdateTaskParams(
	request mcp.CallToolRequest,
) (app.BitbucketUpdateTaskParams, *mcp.CallToolResult) {
	// Extract required parameters
	prID, err := request.RequireInt("pr_id")
	if err != nil {
		return app.BitbucketUpdateTaskParams{}, mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err)
	}

	taskID, err := request.RequireInt("task_id")
	if err != nil {
		return app.BitbucketUpdateTaskParams{}, mcp.NewToolResultErrorFromErr("Missing or invalid task_id parameter", err)
	}

	repoOwner, err := request.RequireString("repo_owner")
	if err != nil {
		return app.BitbucketUpdateTaskParams{}, mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err)
	}

	repoName, err := request.RequireString("repo_name")
	if err != nil {
		return app.BitbucketUpdateTaskParams{}, mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err)
	}

	// Optional parameters
	content := request.GetString("content", "")
	state := request.GetString("state", "")
	account := request.GetString("account", "")

	// Either content or state must be provided
	if content == "" && state == "" {
		return app.BitbucketUpdateTaskParams{}, mcp.NewToolResultError("Either content or state must be provided")
	}

	// Validate state if provided
	if state != "" && state != TaskStateResolved && state != TaskStateUnresolved {
		return app.BitbucketUpdateTaskParams{}, mcp.NewToolResultError("State must be either RESOLVED or UNRESOLVED")
	}

	// Create and return parameters for the service layer
	return app.BitbucketUpdateTaskParams{
		PullRequestID: prID,
		TaskID:        taskID,
		RepoOwner:     repoOwner,
		RepoName:      repoName,
		AccountName:   account,
		Content:       content,
		State:         state,
	}, nil
}

// formatUpdateTaskResponse creates an appropriate response message based on what was updated.
func (bc *BitbucketController) formatUpdateTaskResponse(
	task *bitbucket.PullRequestCommentTask,
	taskID int,
	content,
	state string,
) string {
	switch {
	case content != "" && state != "":
		return fmt.Sprintf("Updated task #%d content and marked as %s", taskID, task.State)
	case content != "":
		return fmt.Sprintf("Updated task #%d content", taskID)
	default:
		return fmt.Sprintf("Updated task #%d state to %s", taskID, task.State)
	}
}

// newCreatePRTaskServerTool returns a server tool for creating tasks on a pull request.
func (bc *BitbucketController) newCreatePRTaskServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_create_pr_task",
		mcp.WithDescription("Create a task on a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Task content"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithNumber("comment_id",
			mcp.Description("Comment ID to associate with the task (optional)"),
		),
		mcp.WithString("state",
			mcp.Description("Initial state for the task: RESOLVED or UNRESOLVED (optional, defaults to UNRESOLVED)"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := bc.makeCreatePRTaskHandler()

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// makeCreatePRTaskHandler creates a handler function for the create PR task tool.
// This is split out to reduce the overall function length.
func (bc *BitbucketController) makeCreatePRTaskHandler() func(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_create_pr_task request", "params", request.Params)

		// Extract required parameters
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid content parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		account := request.GetString("account", "")
		state := request.GetString("state", "")
		commentIDStr := request.GetString("comment_id", "")

		// Parse comment_id if provided
		var commentID int64
		if commentIDStr != "" {
			// Try to convert the string to int64
			var commentIDFloat float64
			commentIDFloat, parseErr := strconv.ParseFloat(commentIDStr, 64)
			if parseErr != nil {
				return mcp.NewToolResultErrorFromErr("Invalid comment_id parameter", parseErr), nil
			}
			commentID = int64(commentIDFloat)
		}

		// Validate state if provided
		if state != "" && state != TaskStateResolved && state != TaskStateUnresolved {
			return mcp.NewToolResultError("State must be either RESOLVED or UNRESOLVED"), nil
		}

		// Create parameters for the service layer
		params := app.BitbucketCreateTaskParams{
			PullRequestID: prID,
			Content:       content,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
			State:         state,
			CommentID:     commentID,
		}

		// Call the service to create the task
		task, err := bc.bitbucketService.CreateTask(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create task: %w", err)
		}

		// Format response text
		var responseText string
		if task.Comment != nil {
			responseText = fmt.Sprintf("Created task on PR #%d: %s (on comment #%d)",
				prID, task.Content.Raw, task.Comment.ID)
		} else {
			responseText = fmt.Sprintf("Created task on PR #%d: %s", prID, task.Content.Raw)
		}

		return mcp.NewToolResultText(responseText), nil
	}
}

// NewTools returns the tools for this controller.
func (bc *BitbucketController) NewTools() []server.ServerTool {
	return []server.ServerTool{
		bc.newCreatePRServerTool(),
		bc.newReadPRServerTool(),
		bc.newUpdatePRServerTool(),
		bc.newApprovePRServerTool(),
		bc.newMergePRServerTool(),
		bc.newListPRTasksServerTool(),
		bc.newUpdatePRTaskServerTool(),
		bc.newCreatePRTaskServerTool(),
		bc.newGetPRDiffstatServerTool(),
		bc.newGetPRDiffServerTool(),
		bc.newAddPRCommentServerTool(),
		bc.newGetFileContentServerTool(),
		bc.newRequestPRChangesServerTool(),
		bc.newListPRCommentsServerTool(),
	}
}

// newGetPRDiffstatServerTool returns a server tool for getting PR diffstat.
func (bc *BitbucketController) newGetPRDiffstatServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_get_pr_diffstat",
		mcp.WithDescription("Get the diffstat for a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_get_pr_diffstat request", "params", request.Params)

		// Extract required parameters
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}
		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}
		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}
		account := request.GetString("account", "")

		// Build params for service layer
		params := app.BitbucketGetPRDiffStatParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
		}

		// Call the service
		result, err := bc.bitbucketService.GetPRDiffStat(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get diffstat: %w", err)
		}

		// Marshal result to JSON
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal diffstat to JSON: %w", err)
		}

		// Create summary text
		summaryText := fmt.Sprintf("Diffstat for PR #%d: %d files changed", prID, result.Size)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: summaryText,
				},
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

/*
 * newGetPRDiffServerTool returns a server tool for getting PR diff.
 */
func (bc *BitbucketController) newGetPRDiffServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_get_pr_diff",
		mcp.WithDescription("Get the diff for a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
		mcp.WithString("file_paths",
			mcp.Description("List of file paths to filter the diff (optional, multiple comma-separated values are possible)"),
		),
		mcp.WithNumber("context_lines",
			mcp.Description("Number of context lines to include in the diff (optional)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_get_pr_diff request", "params", request.Params)

		// Extract required parameters
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}
		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}
		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}
		account := request.GetString("account", "")

		// Extract optional file_paths and context_lines using idiomatic helpers
		filePathsStr := request.GetString("file_paths", "")
		var filePaths []string
		if filePathsStr != "" {
			parts := strings.Split(filePathsStr, ",")
			for _, s := range parts {
				s = strings.TrimSpace(s)
				if s != "" {
					filePaths = append(filePaths, s)
				}
			}
		}
		var contextLines *int
		if cl := request.GetInt("context_lines", 0); cl != 0 {
			contextLines = &cl
		}

		// Build params for service layer
		params := app.BitbucketGetPRDiffParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
			FilePaths:     filePaths,
			ContextLines:  contextLines,
		}

		// Call the service
		diff, err := bc.bitbucketService.GetPRDiff(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get diff: %w", err)
		}

		// Create summary text
		summaryText := fmt.Sprintf("Diff for PR #%d in %s/%s", prID, repoOwner, repoName)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: summaryText,
				},
				mcp.NewTextContent(diff),
			},
		}, nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newGetFileContentServerTool returns a server tool for getting file content from a Bitbucket repository.
func (bc *BitbucketController) newGetFileContentServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_get_file_content",
		mcp.WithDescription("Get the content of a file in a Bitbucket repository"),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("file_path",
			mcp.Description("Path to the file in the repository"),
			mcp.Required(),
		),
		mcp.WithString("commit_hash",
			mcp.Description("The SHA hash to fetch file content from. Only commit hashes are supported."),
			mcp.Required(),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_get_file_content request", "params", request.Params)

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}
		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}
		filePath := request.GetString("file_path", "")
		if filePath == "" {
			filePath = request.GetString("path", "")
		}
		if filePath == "" {
			return mcp.NewToolResultError(
				"Missing or invalid file_path parameter: required argument \"file_path\" not found",
			), nil
		}
		commitHash, err := request.RequireString("commit_hash")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid commit_hash parameter", err), nil
		}
		account := request.GetString("account", "")

		params := app.BitbucketGetFileContentParams{
			AccountName: account,
			RepoOwner:   repoOwner,
			RepoName:    repoName,
			Commit:      commitHash,
			Path:        filePath,
		}

		result, err := bc.bitbucketService.GetFileContent(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get file content: %w", err)
		}

		// Compose summary
		summaryText := fmt.Sprintf("File content for %s at %s/%s", filePath, repoOwner, repoName)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: summaryText,
				},
				mcp.NewTextContent(result.Content),
			},
		}, nil
	}
	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newAddPRCommentServerTool returns a server tool for adding a comment to a pull request.
func (bc *BitbucketController) newAddPRCommentServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_add_pr_comment",
		mcp.WithDescription("Add a comment to a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("comment_text",
			mcp.Description("The comment content"),
			mcp.Required(),
		),
		mcp.WithString("file_path",
			mcp.Description("Path to the file for inline comments (optional)"),
		),
		mcp.WithNumber("line_number_from",
			mcp.Description("Anchor line in the old version of the file (optional)"),
		),
		mcp.WithNumber("line_number_to",
			mcp.Description("Anchor line in the new version of the file (optional)"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
		mcp.WithBoolean("pending",
			mcp.Description("Create as a pending comment (optional, defaults to false)"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_add_pr_comment request", "params", request.Params)

		// Required parameters
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}
		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}
		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}
		account := request.GetString("account", "")
		commentText, err := request.RequireString("comment_text")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid comment_text parameter", err), nil
		}

		// Optional parameters
		filePath := request.GetString("file_path", "")
		lineFrom := request.GetInt("line_number_from", 0)
		lineTo := request.GetInt("line_number_to", 0)
		pending := request.GetBool("pending", false)

		params := app.BitbucketAddPRCommentParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
			Content:       commentText,
			FilePath:      filePath,
			LineFrom:      lineFrom,
			LineTo:        lineTo,
			Pending:       pending,
		}

		commentID, content, err := bc.bitbucketService.AddPRComment(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to add PR comment: %w", err)
		}

		resultText := fmt.Sprintf("Added comment #%d: %s", commentID, content)
		return mcp.NewToolResultText(resultText), nil
	}
	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newRequestPRChangesServerTool returns a server tool for requesting changes on a pull request.

func (bc *BitbucketController) newRequestPRChangesServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_request_pr_changes",
		mcp.WithDescription("Request changes on a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_request_pr_changes request", "params", request.Params)

		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}
		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}
		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}
		account := request.GetString("account", "")

		params := app.BitbucketRequestPRChangesParams{
			PullRequestID: prID,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
			AccountName:   account,
		}

		status, timestamp, err := bc.bitbucketService.RequestPRChanges(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to request PR changes: %w", err)
		}

		return mcp.NewToolResultText(
			fmt.Sprintf("Requested changes for pull request #%d: %s at %v", prID, status, timestamp),
		), nil
	}
	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// newListPRCommentsServerTool returns a server tool for listing comments on a pull request.
func (bc *BitbucketController) newListPRCommentsServerTool() server.ServerTool {
	tool := mcp.NewTool(
		"bitbucket_list_pr_comments",
		mcp.WithDescription("List all comments on a pull request in Bitbucket"),
		mcp.WithNumber("pr_id",
			mcp.Description("Pull request ID"),
			mcp.Required(),
		),
		mcp.WithString("repo_owner",
			mcp.Description("Repository owner (username/workspace)"),
			mcp.Required(),
		),
		mcp.WithString("repo_name",
			mcp.Description("Repository name (slug)"),
			mcp.Required(),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bc.logger.Debug("Received bitbucket_list_pr_comments request", "params", request.Params)

		// Extract required parameters using RequireXXX methods
		prID, err := request.RequireInt("pr_id")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid pr_id parameter", err), nil
		}

		repoOwner, err := request.RequireString("repo_owner")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_owner parameter", err), nil
		}

		repoName, err := request.RequireString("repo_name")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Missing or invalid repo_name parameter", err), nil
		}

		// Optional parameters
		account := request.GetString("account", "")

		// Create parameters for the service layer
		params := app.BitbucketListPRCommentsParams{
			PullRequestID: prID,
			AccountName:   account,
			RepoOwner:     repoOwner,
			RepoName:      repoName,
		}

		// Call the service to list PR comments
		comments, err := bc.bitbucketService.ListPRComments(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list PR comments: %w", err)
		}

		// Convert comments to JSON for the resource
		commentsJSON, err := json.MarshalIndent(comments, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal PR comments to JSON: %w", err)
		}

		// Create a summary text for the comments
		summaryText := fmt.Sprintf("Found %d comments on pull request #%d", len(comments.Values), prID)

		// Return both a summary text and the full comments data as a resource
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: summaryText,
				},

				// Sending json as text since some clients (Cursor)
				// do not support resources (at least not yet)
				mcp.NewTextContent(string(commentsJSON)),
			},
		}, nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}
