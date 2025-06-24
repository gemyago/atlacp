package controllers

import (
	"context"
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
		mcp.WithString("description",
			mcp.Description("Pull request description"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
		mcp.WithArray("reviewers",
			mcp.Description("Usernames of reviewers to assign"),
		),
	)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// TODO: Implement handler for create PR
		return mcp.NewToolResultText("CreatePR functionality not implemented yet"), nil
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
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// TODO: Implement handler for read PR
		return mcp.NewToolResultText("ReadPR functionality not implemented yet"), nil
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

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// TODO: Implement handler for update PR
		return mcp.NewToolResultText("UpdatePR functionality not implemented yet"), nil
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
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// TODO: Implement handler for approve PR
		return mcp.NewToolResultText("ApprovePR functionality not implemented yet"), nil
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
		mcp.WithString("merge_strategy",
			mcp.Description("Merge strategy to use"),
			mcp.Enum("merge_commit", "squash", "fast_forward"),
		),
		mcp.WithString("commit_message",
			mcp.Description("Custom commit message for merge commit"),
		),
		mcp.WithString("account",
			mcp.Description("Atlassian account name to use (optional, uses default if not specified)"),
		),
	)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// TODO: Implement handler for merge PR
		return mcp.NewToolResultText("MergePR functionality not implemented yet"), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// NewTools returns all Bitbucket tools.
// Satisfies the ToolsFactory interface.
func (bc *BitbucketController) NewTools() []server.ServerTool {
	return []server.ServerTool{
		bc.newCreatePRServerTool(),
		bc.newReadPRServerTool(),
		bc.newUpdatePRServerTool(),
		bc.newApprovePRServerTool(),
		bc.newMergePRServerTool(),
	}
}
