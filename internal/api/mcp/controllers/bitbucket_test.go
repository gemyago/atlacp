package controllers

import (
	"fmt"
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

			// Create test data with randomized values using faker
			title := "PR-" + faker.Sentence()
			sourceBranch := "feature/" + faker.Username()
			targetBranch := "main"
			description := faker.Paragraph()
			reviewers := []string{"user-" + faker.Username(), "user-" + faker.Username()}
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()

			// Create expected parameters and response
			expectedParams := app.BitbucketCreatePRParams{
				Title:        title,
				SourceBranch: sourceBranch,
				DestBranch:   targetBranch,
				Description:  description,
				Reviewers:    reviewers,
				RepoOwner:    repoOwner,
				RepoName:     repoName,
			}

			createdOn := time.Now()
			updatedOn := time.Now()
			prID := int(faker.RandomUnixTime()) % 1000000 // Generate a random PR ID

			// Create expected PR response with random values
			expectedPR := &bitbucket.PullRequest{
				ID:          prID,
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
					DisplayName: "User-" + faker.Name(),
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
						"repo_owner":    repoOwner,
						"repo_name":     repoName,
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
			assert.Contains(t, content.Text, fmt.Sprintf("Created pull request #%d", prID))
			assert.Contains(t, content.Text, title)
		})

		t.Run("should handle missing required parameters", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create request missing required parameters but with other valid data
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_create_pr",
					Arguments: map[string]interface{}{
						// Missing title
						"source_branch": "feature/" + faker.Username(),
						"target_branch": "main",
						"repo_owner":    repoOwner,
						"repo_name":     repoName,
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

		t.Run("should handle ReadPR call successfully", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create test data with randomized values using faker
			prID := int(faker.RandomUnixTime()) % 1000000 // Generate a random PR ID within a reasonable range
			accountName := "account-" + faker.Username()
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()

			// Create expected parameters and response
			expectedParams := app.BitbucketReadPRParams{
				PullRequestID: prID,
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
			}

			createdOn := time.Now()
			updatedOn := time.Now()

			// Create expected PR response with random values
			expectedPR := &bitbucket.PullRequest{
				ID:          prID,
				Title:       "PR-" + faker.Sentence(),
				Description: faker.Paragraph(),
				State:       "OPEN",
				Source: bitbucket.PullRequestSource{
					Branch: bitbucket.PullRequestBranch{
						Name: "feature/" + faker.Username(),
					},
				},
				Destination: &bitbucket.PullRequestDestination{
					Branch: bitbucket.PullRequestBranch{
						Name: "main",
					},
				},
				Author: &bitbucket.PullRequestAuthor{
					DisplayName: "User-" + faker.Name(),
				},
				CreatedOn: &createdOn,
				UpdatedOn: &updatedOn,
			}

			// Setup mock expectations with any matcher for context
			mockService.EXPECT().
				ReadPR(mock.Anything, mock.MatchedBy(func(params app.BitbucketReadPRParams) bool {
					return params.PullRequestID == expectedParams.PullRequestID &&
						params.AccountName == expectedParams.AccountName &&
						params.RepoOwner == expectedParams.RepoOwner &&
						params.RepoName == expectedParams.RepoName
				})).
				Return(expectedPR, nil)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_read_pr",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"account":    accountName,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
					},
				},
			}

			// Get the handler
			serverTool := controller.newReadPRServerTool()
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
			assert.Contains(t, content.Text, fmt.Sprintf("Pull request #%d", prID))
			assert.Contains(t, content.Text, expectedPR.Title)
			assert.Contains(t, content.Text, expectedPR.State)
		})

		t.Run("should handle missing PR ID parameter in ReadPR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create request missing required parameters but with other valid data
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			accountName := "account-" + faker.Username()

			// Create request missing required parameters
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_read_pr",
					Arguments: map[string]interface{}{
						// Missing pr_id
						"account":    accountName,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
					},
				},
			}

			// Get the handler
			serverTool := controller.newReadPRServerTool()
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
			assert.Contains(t, content.Text, "Missing or invalid pr_id parameter")
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
