package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/dig"
)

// bitbucketService defines the operations required by the BitbucketController.
// This interface matches the methods from app.BitbucketService that are used by the controller.
type bitbucketService interface {
	CreatePR(ctx context.Context, params app.BitbucketCreatePRParams) (*bitbucket.PullRequest, error)
	ReadPR(ctx context.Context, params app.BitbucketReadPRParams) (*bitbucket.PullRequest, error)
	UpdatePR(ctx context.Context, params app.BitbucketUpdatePRParams) (*bitbucket.PullRequest, error)
	ApprovePR(ctx context.Context, params app.BitbucketApprovePRParams) (*bitbucket.Participant, error)
	MergePR(ctx context.Context, params app.BitbucketMergePRParams) (*bitbucket.PullRequest, error)
}

// Ensure that app.BitbucketService implements bitbucketService.
var _ bitbucketService = (*app.BitbucketService)(nil)

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

		// Create parameters for the service layer
		params := app.BitbucketCreatePRParams{
			Title:        title,
			SourceBranch: sourceBranch,
			DestBranch:   targetBranch,
			Description:  description,
			AccountName:  account,
			RepoOwner:    repoOwner,
			RepoName:     repoName,
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

		// At least one of title or description must be provided
		title := request.GetString("title", "")
		description := request.GetString("description", "")

		if title == "" && description == "" {
			return mcp.NewToolResultError("At least one of title or description must be provided"), nil
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
		closeSourceBranchStr := request.GetString("close_source_branch", "")
		account := request.GetString("account", "")

		// Parse boolean parameter
		var closeSourceBranch bool
		if closeSourceBranchStr != "" {
			closeSourceBranch = closeSourceBranchStr == "true"
		}

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

// NewTools returns the tools for this controller.
func (bc *BitbucketController) NewTools() []server.ServerTool {
	return []server.ServerTool{
		bc.newCreatePRServerTool(),
		bc.newReadPRServerTool(),
		bc.newUpdatePRServerTool(),
		bc.newApprovePRServerTool(),
		bc.newMergePRServerTool(),
	}
}
