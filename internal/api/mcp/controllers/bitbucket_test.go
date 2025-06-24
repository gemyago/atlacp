package controllers

import (
	"testing"
	"time"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/gemyago/atlacp/internal/testing/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// The testing.T parameter is used for logging in test context.
func makeBitbucketControllerDeps(t *testing.T) BitbucketControllerDeps {
	// Create a mock bitbucketService for testing
	mockBitbucketService := NewMockbitbucketService(t)

	// Use t for test logging context
	logger := diag.RootTestLogger().With("test", t.Name())

	return BitbucketControllerDeps{
		RootLogger:       logger,
		BitbucketService: mockBitbucketService,
	}
}

func TestBitbucketController(t *testing.T) {
	t.Run("should create Bitbucket controller with dependencies", func(t *testing.T) {
		deps := makeBitbucketControllerDeps(t)

		controller := NewBitbucketController(deps)

		require.NotNil(t, controller)
		require.NotNil(t, controller.logger)
		require.NotNil(t, controller.bitbucketService)
	})

	t.Run("tool definitions", func(t *testing.T) {
		t.Run("should define CreatePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newCreatePRServerTool()

			assert.Equal(t, "bitbucket_create_pr", serverTool.Tool.Name)
			assert.Equal(t, "Create a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ReadPR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newReadPRServerTool()

			assert.Equal(t, "bitbucket_read_pr", serverTool.Tool.Name)
			assert.Equal(t, "Get pull request details from Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define UpdatePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newUpdatePRServerTool()

			assert.Equal(t, "bitbucket_update_pr", serverTool.Tool.Name)
			assert.Equal(t, "Update a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ApprovePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newApprovePRServerTool()

			assert.Equal(t, "bitbucket_approve_pr", serverTool.Tool.Name)
			assert.Equal(t, "Approve a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define MergePR tool correctly", func(t *testing.T) {
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newMergePRServerTool()

			assert.Equal(t, "bitbucket_merge_pr", serverTool.Tool.Name)
			assert.Equal(t, "Merge a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})
	})

	t.Run("should register all tools", func(t *testing.T) {
		deps := makeBitbucketControllerDeps(t)
		controller := NewBitbucketController(deps)

		tools := controller.NewTools()

		require.Len(t, tools, 5) // Expect 5 tools: create, read, update, approve, merge
		toolNames := make([]string, len(tools))
		for i, tool := range tools {
			toolNames[i] = tool.Tool.Name
		}
		assert.Contains(t, toolNames, "bitbucket_create_pr")
		assert.Contains(t, toolNames, "bitbucket_read_pr")
		assert.Contains(t, toolNames, "bitbucket_update_pr")
		assert.Contains(t, toolNames, "bitbucket_approve_pr")
		assert.Contains(t, toolNames, "bitbucket_merge_pr")
	})

	t.Run("handler implementations", func(t *testing.T) {
		t.Run("should handle CreatePR call successfully", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create test data
			title := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Username()
			targetBranch := "main"
			description := faker.Paragraph()
			reviewers := []string{faker.Username(), faker.Username()}

			// Create expected parameters and response
			expectedParams := app.BitbucketCreatePRParams{
				Title:        title,
				SourceBranch: sourceBranch,
				DestBranch:   targetBranch,
				Description:  description,
				Reviewers:    reviewers,
				RepoOwner:    "your-workspace",  // Hardcoded value in the implementation
				RepoName:     "your-repository", // Hardcoded value in the implementation
			}

			createdOn := time.Now()
			updatedOn := time.Now()

			// Create expected PR response
			expectedPR := &bitbucket.PullRequest{
				ID:          12345,
				Title:       title,
				Description: description,
				State:       "OPEN",
				Source: bitbucket.PullRequestSource{
					Branch: bitbucket.PullRequestBranch{
						Name: sourceBranch,
					},
				},
				Destination: &bitbucket.PullRequestDestination{
					Branch: bitbucket.PullRequestBranch{
						Name: targetBranch,
					},
				},
				Author: &bitbucket.PullRequestAuthor{
					DisplayName: "Test User",
				},
				CreatedOn: &createdOn,
				UpdatedOn: &updatedOn,
			}

			// Setup mock expectations with any matcher for context
			mockService.EXPECT().
				CreatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketCreatePRParams) bool {
					return params.Title == expectedParams.Title &&
						params.SourceBranch == expectedParams.SourceBranch &&
						params.DestBranch == expectedParams.DestBranch &&
						params.Description == expectedParams.Description &&
						params.RepoOwner == expectedParams.RepoOwner &&
						params.RepoName == expectedParams.RepoName
				})).
				Return(expectedPR, nil)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_create_pr",
					Arguments: map[string]interface{}{
						"title":         title,
						"source_branch": sourceBranch,
						"target_branch": targetBranch,
						"description":   description,
						"reviewers":     reviewers,
					},
				},
			}

			// Get the handler
			serverTool := controller.newCreatePRServerTool()
			handler := serverTool.Handler

			// Act
			result, err := handler(ctx, request)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.False(t, result.IsError)

			// Verify the content of the result
			content, ok := result.Content[0].(mcp.TextContent)
			require.True(t, ok, "Result content should be text content")
			assert.Contains(t, content.Text, "Created pull request #12345")
			assert.Contains(t, content.Text, title)
		})

		t.Run("should handle missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create request missing required parameters
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_create_pr",
					Arguments: map[string]interface{}{
						// Missing title
						"source_branch": "feature/branch",
						"target_branch": "main",
					},
				},
			}

			// Get the handler
			serverTool := controller.newCreatePRServerTool()
			handler := serverTool.Handler

			// Act
			result, err := handler(ctx, request)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.True(t, result.IsError)

			// Verify error message
			content, ok := result.Content[0].(mcp.TextContent)
			require.True(t, ok, "Error content should be text content")
			assert.Contains(t, content.Text, "Missing or invalid title parameter")
		})
	})

	// Example of how to use mocks.GetMock for retrieving the mock instance:
	// t.Run("should handle CreatePR call", func(t *testing.T) {
	//     deps := makeBitbucketControllerDeps(t)
	//     mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
	//     controller := NewBitbucketController(deps)
	//
	//     // Setup mock expectations
	//     mockService.EXPECT().CreatePR(...).Return(...)
	//
	//     // Test the handler
	//     // ...
	// })

	// Future tests for each handler implementation will go here
	// These will be added once we implement the actual handlers
}
