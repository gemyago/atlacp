package controllers

import (
	"errors"
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

		t.Run("should handle UpdatePR call successfully", func(t *testing.T) {
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
			newTitle := "Updated-PR-" + faker.Sentence()
			newDescription := faker.Paragraph()

			// Create expected parameters and response
			expectedParams := app.BitbucketUpdatePRParams{
				PullRequestID: prID,
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
				Title:         newTitle,
				Description:   newDescription,
			}

			createdOn := time.Now()
			updatedOn := time.Now()

			// Create expected PR response with random values
			expectedPR := &bitbucket.PullRequest{
				ID:          prID,
				Title:       newTitle,
				Description: newDescription,
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
				UpdatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketUpdatePRParams) bool {
					return params.PullRequestID == expectedParams.PullRequestID &&
						params.AccountName == expectedParams.AccountName &&
						params.RepoOwner == expectedParams.RepoOwner &&
						params.RepoName == expectedParams.RepoName &&
						params.Title == expectedParams.Title &&
						params.Description == expectedParams.Description
				})).
				Return(expectedPR, nil)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_update_pr",
					Arguments: map[string]interface{}{
						"pr_id":       prID,
						"account":     accountName,
						"repo_owner":  repoOwner,
						"repo_name":   repoName,
						"title":       newTitle,
						"description": newDescription,
					},
				},
			}

			// Get the handler
			serverTool := controller.newUpdatePRServerTool()
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
			assert.Contains(t, content.Text, fmt.Sprintf("Updated pull request #%d", prID))
			assert.Contains(t, content.Text, newTitle)
		})

		t.Run("should handle missing required parameters in UpdatePR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create request missing required parameters but with other valid data
			accountName := "account-" + faker.Username()
			repoName := "repo-" + faker.Word()

			// Create request missing PR ID and repo_owner
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_update_pr",
					Arguments: map[string]interface{}{
						// Missing pr_id and repo_owner
						"account":     accountName,
						"repo_name":   repoName,
						"title":       "Updated title",
						"description": faker.Paragraph(),
					},
				},
			}

			// Get the handler
			serverTool := controller.newUpdatePRServerTool()
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

		t.Run("should handle missing both title and description in UpdatePR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create test data with randomized values
			prID := int(faker.RandomUnixTime()) % 1000000
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()

			// Create request with required fields but missing both title and description
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_update_pr",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
						// Missing both title and description
					},
				},
			}

			// Get the handler
			serverTool := controller.newUpdatePRServerTool()
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
			assert.Contains(t, content.Text, "At least one of title or description must be provided")
		})

		t.Run("should handle ApprovePR call successfully", func(t *testing.T) {
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

			// Generate random user data for returned participant
			displayName := "User-" + faker.Name()
			username := "user-" + faker.Username()
			role := "REVIEWER" // This is typically a fixed value in Bitbucket

			// Create expected parameters and response
			expectedParams := app.BitbucketApprovePRParams{
				PullRequestID: prID,
				AccountName:   accountName,
				RepoOwner:     repoOwner,
				RepoName:      repoName,
			}

			// Create expected PR participant response
			expectedParticipant := &bitbucket.Participant{
				User: bitbucket.PullRequestAuthor{
					DisplayName: displayName,
					Username:    username,
				},
				Role:     role,
				Approved: true,
				State:    "APPROVED",
				Type:     "participant",
			}

			// Setup mock expectations with any matcher for context
			mockService.EXPECT().
				ApprovePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketApprovePRParams) bool {
					return params.PullRequestID == expectedParams.PullRequestID &&
						params.AccountName == expectedParams.AccountName &&
						params.RepoOwner == expectedParams.RepoOwner &&
						params.RepoName == expectedParams.RepoName
				})).
				Return(expectedParticipant, nil)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_approve_pr",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"account":    accountName,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
					},
				},
			}

			// Get the handler
			serverTool := controller.newApprovePRServerTool()
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
			assert.Contains(t, content.Text, fmt.Sprintf("Pull request #%d approved", prID))
			assert.Contains(t, content.Text, displayName)
		})

		t.Run("should handle missing required parameters in ApprovePR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create request missing required parameters but with other valid data
			accountName := "account-" + faker.Username()
			repoOwner := "workspace-" + faker.Username()

			// Create request missing PR ID
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_approve_pr",
					Arguments: map[string]interface{}{
						// Missing pr_id and repo_name
						"account":    accountName,
						"repo_owner": repoOwner,
					},
				},
			}

			// Get the handler
			serverTool := controller.newApprovePRServerTool()
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

		t.Run("should handle service error in ApprovePR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create test data with randomized values using faker
			prID := int(faker.RandomUnixTime()) % 1000000
			accountName := "account-" + faker.Username()
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()

			// Create expected error
			expectedError := errors.New("failed to approve pull request: " + faker.Sentence())

			// Setup mock to return an error
			mockService.EXPECT().
				ApprovePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketApprovePRParams) bool {
					return params.PullRequestID == prID &&
						params.AccountName == accountName &&
						params.RepoOwner == repoOwner &&
						params.RepoName == repoName
				})).
				Return(nil, expectedError)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_approve_pr",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"account":    accountName,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
					},
				},
			}

			// Get the handler
			serverTool := controller.newApprovePRServerTool()
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
			assert.Contains(t, content.Text, expectedError.Error())
		})

		t.Run("should handle MergePR call successfully", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create test data with randomized values using faker
			prID := int(faker.RandomUnixTime()) % 1000000
			accountName := "account-" + faker.Username()
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			mergeStrategy := "squash"
			commitMessage := "Merge PR: " + faker.Sentence()
			closeSourceBranch := true

			// Create expected parameters
			expectedParams := app.BitbucketMergePRParams{
				PullRequestID:     prID,
				AccountName:       accountName,
				RepoOwner:         repoOwner,
				RepoName:          repoName,
				MergeStrategy:     mergeStrategy,
				Message:           commitMessage,
				CloseSourceBranch: closeSourceBranch,
			}

			// Create expected PR response
			expectedPR := &bitbucket.PullRequest{
				ID:          prID,
				Title:       "PR-" + faker.Sentence(),
				Description: faker.Paragraph(),
				State:       "MERGED",
			}

			// Setup mock expectations
			mockService.EXPECT().
				MergePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketMergePRParams) bool {
					return params.PullRequestID == expectedParams.PullRequestID &&
						params.AccountName == expectedParams.AccountName &&
						params.RepoOwner == expectedParams.RepoOwner &&
						params.RepoName == expectedParams.RepoName &&
						params.MergeStrategy == expectedParams.MergeStrategy &&
						params.Message == expectedParams.Message &&
						params.CloseSourceBranch == expectedParams.CloseSourceBranch
				})).
				Return(expectedPR, nil)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_merge_pr",
					Arguments: map[string]interface{}{
						"pr_id":               prID,
						"account":             accountName,
						"repo_owner":          repoOwner,
						"repo_name":           repoName,
						"merge_strategy":      mergeStrategy,
						"commit_message":      commitMessage,
						"close_source_branch": "true",
					},
				},
			}

			// Get the handler
			serverTool := controller.newMergePRServerTool()
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
			assert.Contains(t, content.Text, fmt.Sprintf("Pull request #%d successfully merged", prID))
			assert.Contains(t, content.Text, "using squash strategy")
			assert.Contains(t, content.Text, "source branch was closed")
		})

		t.Run("should handle missing required parameters in MergePR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create request missing required parameters
			accountName := "account-" + faker.Username()

			// Missing repo_owner and repo_name
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_merge_pr",
					Arguments: map[string]interface{}{
						"pr_id":   123,
						"account": accountName,
						// Missing repo_owner and repo_name
					},
				},
			}

			// Get the handler
			serverTool := controller.newMergePRServerTool()
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
			assert.Contains(t, content.Text, "Missing or invalid repo_owner parameter")
		})

		t.Run("should handle service error in MergePR", func(t *testing.T) {
			// Arrange
			deps := makeBitbucketControllerDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Create test data with randomized values
			prID := int(faker.RandomUnixTime()) % 1000000
			accountName := "account-" + faker.Username()
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()

			// Create expected error
			expectedError := errors.New("failed to merge pull request: " + faker.Sentence())

			// Setup mock to return an error
			mockService.EXPECT().
				MergePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketMergePRParams) bool {
					return params.PullRequestID == prID &&
						params.AccountName == accountName &&
						params.RepoOwner == repoOwner &&
						params.RepoName == repoName
				})).
				Return(nil, expectedError)

			// Create the request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_merge_pr",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"account":    accountName,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
					},
				},
			}

			// Get the handler
			serverTool := controller.newMergePRServerTool()
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
			assert.Contains(t, content.Text, expectedError.Error())
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
