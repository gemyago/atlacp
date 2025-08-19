package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/gemyago/atlacp/internal/app"
	"github.com/gemyago/atlacp/internal/diag"
	"github.com/gemyago/atlacp/internal/services/bitbucket"
	"github.com/gemyago/atlacp/internal/testing/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBitbucketController(t *testing.T) {
	makeMockDeps := func(t *testing.T) BitbucketControllerDeps {
		// Create a mock bitbucketService for testing
		mockBitbucketService := NewMockbitbucketService(t)

		// Use t for test logging context
		logger := diag.RootTestLogger().With("test", t.Name())

		return BitbucketControllerDeps{
			RootLogger:       logger,
			BitbucketService: mockBitbucketService,
		}
	}

	t.Run("should create Bitbucket controller with dependencies", func(t *testing.T) {
		deps := makeMockDeps(t)

		controller := NewBitbucketController(deps)

		require.NotNil(t, controller)
		require.NotNil(t, controller.logger)
		require.NotNil(t, controller.bitbucketService)
	})

	t.Run("tool definitions", func(t *testing.T) {
		t.Run("should define CreatePR tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newCreatePRServerTool()

			assert.Equal(t, "bitbucket_create_pr", serverTool.Tool.Name)
			assert.Equal(t, "Create a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ReadPR tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newReadPRServerTool()

			assert.Equal(t, "bitbucket_read_pr", serverTool.Tool.Name)
			assert.Equal(t, "Get pull request details from Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define UpdatePR tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newUpdatePRServerTool()

			assert.Equal(t, "bitbucket_update_pr", serverTool.Tool.Name)
			assert.Equal(t, "Update a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ApprovePR tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newApprovePRServerTool()

			assert.Equal(t, "bitbucket_approve_pr", serverTool.Tool.Name)
			assert.Equal(t, "Approve a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define MergePR tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newMergePRServerTool()

			assert.Equal(t, "bitbucket_merge_pr", serverTool.Tool.Name)
			assert.Equal(t, "Merge a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define ListPRTasks tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newListPRTasksServerTool()

			assert.Equal(t, "bitbucket_list_pr_tasks", serverTool.Tool.Name)
			assert.Equal(t, "List tasks on a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define UpdatePRTask tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newUpdatePRTaskServerTool()

			assert.Equal(t, "bitbucket_update_pr_task", serverTool.Tool.Name)
			assert.Equal(t, "Update a task on a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})

		t.Run("should define CreatePRTask tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newCreatePRTaskServerTool()

			assert.Equal(t, "bitbucket_create_pr_task", serverTool.Tool.Name)
			assert.Equal(t, "Create a task on a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})
		t.Run("should define GetPRDiffstat tool correctly", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)

			serverTool := controller.newGetPRDiffstatServerTool()

			assert.Equal(t, "bitbucket_get_pr_diffstat", serverTool.Tool.Name)
			assert.Equal(t, "Get the diffstat for a pull request in Bitbucket", serverTool.Tool.Description)
			assert.NotNil(t, serverTool.Tool.InputSchema)
			assert.NotNil(t, serverTool.Handler)
		})
	})

	t.Run("should register all tools", func(t *testing.T) {
		deps := makeMockDeps(t)
		controller := NewBitbucketController(deps)

		tools := controller.NewTools()

		require.Len(t, tools, 9) // 9 tools: create, read, update, approve, merge, list, update, create task, get diffstat
		toolNames := make([]string, len(tools))
		for i, tool := range tools {
			toolNames[i] = tool.Tool.Name
		}
		assert.Contains(t, toolNames, "bitbucket_create_pr")
		assert.Contains(t, toolNames, "bitbucket_read_pr")
		assert.Contains(t, toolNames, "bitbucket_update_pr")
		assert.Contains(t, toolNames, "bitbucket_approve_pr")
		assert.Contains(t, toolNames, "bitbucket_merge_pr")
		assert.Contains(t, toolNames, "bitbucket_list_pr_tasks")
		assert.Contains(t, toolNames, "bitbucket_update_pr_task")
		assert.Contains(t, toolNames, "bitbucket_create_pr_task")
		assert.Contains(t, toolNames, "bitbucket_get_pr_diffstat")
	})

	t.Run("handlers", func(t *testing.T) {
		t.Run("bitbucket_create_pr", func(t *testing.T) {
			t.Run("should handle CreatePR call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				title := "PR-" + faker.Sentence()
				sourceBranch := "feature/" + faker.Username()
				targetBranch := "main"
				description := faker.Paragraph()
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()

				// Create expected parameters and response
				expectedParams := app.BitbucketCreatePRParams{
					Title:        title,
					SourceBranch: sourceBranch,
					DestBranch:   targetBranch,
					Description:  description,
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

			t.Run("should handle draft pull request creation", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				title := "DRAFT: " + faker.Sentence()
				sourceBranch := "feature/" + faker.Username()
				targetBranch := "main"
				description := faker.Paragraph()
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()

				// Create expected parameters and response
				expectedParams := app.BitbucketCreatePRParams{
					Title:        title,
					SourceBranch: sourceBranch,
					DestBranch:   targetBranch,
					Description:  description,
					RepoOwner:    repoOwner,
					RepoName:     repoName,
					Draft:        lo.ToPtr(true), // Expect draft to be true
				}

				prID := int(faker.RandomUnixTime()) % 1000000 // Generate a random PR ID

				// Create expected PR response with random values and Draft=true
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
					Draft: lo.ToPtr(true), // PR is a draft
				}

				// Setup mock expectations with draft=true
				mockService.EXPECT().
					CreatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketCreatePRParams) bool {
						return params.Title == expectedParams.Title &&
							params.SourceBranch == expectedParams.SourceBranch &&
							params.DestBranch == expectedParams.DestBranch &&
							params.Description == expectedParams.Description &&
							params.RepoOwner == expectedParams.RepoOwner &&
							params.RepoName == expectedParams.RepoName &&
							*params.Draft == true // Verify draft flag is set to true
					})).
					Return(expectedPR, nil)

				// Create the request with draft=true
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr",
						Arguments: map[string]interface{}{
							"title":         title,
							"source_branch": sourceBranch,
							"target_branch": targetBranch,
							"description":   description,
							"repo_owner":    repoOwner,
							"repo_name":     repoName,
							"draft":         "true", // Set as draft PR
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
				deps := makeMockDeps(t)
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

			t.Run("should handle missing source_branch parameter", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr",
						Arguments: map[string]interface{}{
							"title": "PR-" + faker.Sentence(),
							// Missing source_branch
							"target_branch": "main",
							"repo_owner":    "workspace-" + faker.Username(),
							"repo_name":     "repo-" + faker.Word(),
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
				assert.Contains(t, content.Text, "Missing or invalid source_branch parameter")
			})

			t.Run("should handle missing target_branch parameter", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr",
						Arguments: map[string]interface{}{
							"title":         "PR-" + faker.Sentence(),
							"source_branch": "feature/" + faker.Username(),
							// Missing target_branch
							"repo_owner": "workspace-" + faker.Username(),
							"repo_name":  "repo-" + faker.Word(),
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
				assert.Contains(t, content.Text, "Missing or invalid target_branch parameter")
			})

			t.Run("should handle missing repo_owner parameter", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr",
						Arguments: map[string]interface{}{
							"title":         "PR-" + faker.Sentence(),
							"source_branch": "feature/" + faker.Username(),
							"target_branch": "main",
							// Missing repo_owner
							"repo_name": "repo-" + faker.Word(),
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
				assert.Contains(t, content.Text, "Missing or invalid repo_owner parameter")
			})

			t.Run("should handle missing repo_name parameter", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr",
						Arguments: map[string]interface{}{
							"title":         "PR-" + faker.Sentence(),
							"source_branch": "feature/" + faker.Username(),
							"target_branch": "main",
							"repo_owner":    "workspace-" + faker.Username(),
							// Missing repo_name
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
				assert.Contains(t, content.Text, "Missing or invalid repo_name parameter")
			})

			t.Run("should handle service error in CreatePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				title := "PR-" + faker.Sentence()
				sourceBranch := "feature/" + faker.Username()
				targetBranch := "main"
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()

				// Create expected error
				expectedError := errors.New("failed to create pull request: " + faker.Sentence())

				// Setup mock to return an error
				mockService.EXPECT().
					CreatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketCreatePRParams) bool {
						return params.Title == title &&
							params.SourceBranch == sourceBranch &&
							params.DestBranch == targetBranch &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName
					})).
					Return(nil, expectedError)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr",
						Arguments: map[string]interface{}{
							"title":         title,
							"source_branch": sourceBranch,
							"target_branch": targetBranch,
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
				require.Error(t, err)
				assert.Contains(t, err.Error(), expectedError.Error())
				assert.Nil(t, result)
			})
		})

		t.Run("bitbucket_read_pr", func(t *testing.T) {
			t.Run("should handle ReadPR call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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

				// Verify the result has both text content items
				require.Len(t, result.Content, 2, "Result should have two text content items")

				// Check first text content (summary)
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok, "First content item should be text")
				assert.Contains(t, textContent.Text, fmt.Sprintf("Pull request #%d", prID))

				// Check second text content (JSON data)
				jsonContent, ok := result.Content[1].(mcp.TextContent)
				require.True(t, ok, "Second content item should also be text")

				// Parse the JSON back to a PullRequest struct
				var receivedPR bitbucket.PullRequest
				err = json.Unmarshal([]byte(jsonContent.Text), &receivedPR)
				require.NoError(t, err, "Should be able to parse JSON back to PullRequest struct")

				// Compare structs directly
				assert.Equal(t, expectedPR.Title, receivedPR.Title)
				assert.Equal(t, expectedPR.State, receivedPR.State)
			})

			t.Run("should return PR details as embedded resource", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				accountName := "account-" + faker.Username()

				// Use testing utility to generate random pull request
				expectedPR := bitbucket.NewRandomPullRequest(
					bitbucket.WithPullRequestID(prID),
				)
				expectedPR.MergeCommit = &bitbucket.PullRequestCommit{
					Hash: faker.UUIDHyphenated(),
				}

				// Setup mock expectations
				mockService.EXPECT().
					ReadPR(mock.Anything, mock.MatchedBy(func(params app.BitbucketReadPRParams) bool {
						return params.PullRequestID == prID &&
							params.AccountName == accountName &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName
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

				// Verify the result has two text content items instead of text + resource
				require.Len(t, result.Content, 2, "Result should have two text content items")

				// Check first text content (summary)
				textContent, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok, "First content item should be text")
				assert.Contains(t, textContent.Text, fmt.Sprintf("Pull request #%d", prID))

				// Check second text content (JSON data)
				jsonContent, ok := result.Content[1].(mcp.TextContent)
				require.True(t, ok, "Second content item should also be text")

				// Parse the JSON back to a PullRequest struct
				var receivedPR bitbucket.PullRequest
				err = json.Unmarshal([]byte(jsonContent.Text), &receivedPR)
				require.NoError(t, err, "Should be able to parse JSON back to PullRequest struct")

				// Compare structs directly
				assert.Equal(t, *expectedPR, receivedPR)
			})

			t.Run("should handle missing PR ID parameter in ReadPR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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

			t.Run("should handle missing repo_owner parameter in ReadPR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_read_pr",
						Arguments: map[string]interface{}{
							"pr_id":   int(faker.RandomUnixTime()) % 1000000,
							"account": "account-" + faker.Username(),
							// Missing repo_owner
							"repo_name": "repo-" + faker.Word(),
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
				assert.Contains(t, content.Text, "Missing or invalid repo_owner parameter")
			})

			t.Run("should handle missing repo_name parameter in ReadPR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_read_pr",
						Arguments: map[string]interface{}{
							"pr_id":      int(faker.RandomUnixTime()) % 1000000,
							"account":    "account-" + faker.Username(),
							"repo_owner": "workspace-" + faker.Username(),
							// Missing repo_name
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
				assert.Contains(t, content.Text, "Missing or invalid repo_name parameter")
			})

			t.Run("should handle service error in ReadPR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := int(faker.RandomUnixTime()) % 1000000
				accountName := "account-" + faker.Username()
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()

				// Create expected error
				expectedError := errors.New("failed to read pull request: " + faker.Sentence())

				// Setup mock to return an error
				mockService.EXPECT().
					ReadPR(mock.Anything, mock.MatchedBy(func(params app.BitbucketReadPRParams) bool {
						return params.PullRequestID == prID &&
							params.AccountName == accountName &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName
					})).
					Return(nil, expectedError)

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
				require.Error(t, err)
				assert.Contains(t, err.Error(), expectedError.Error())
				assert.Nil(t, result)
			})
		})

		t.Run("bitbucket_update_pr", func(t *testing.T) {
			t.Run("should handle UpdatePR call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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
				deps := makeMockDeps(t)
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

			t.Run("should handle missing repo_owner parameter in UpdatePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr",
						Arguments: map[string]interface{}{
							"pr_id":   int(faker.RandomUnixTime()) % 1000000,
							"account": "account-" + faker.Username(),
							// Missing repo_owner
							"repo_name":   "repo-" + faker.Word(),
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
				assert.Contains(t, content.Text, "Missing or invalid repo_owner parameter")
			})

			t.Run("should handle missing repo_name parameter in UpdatePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr",
						Arguments: map[string]interface{}{
							"pr_id":      int(faker.RandomUnixTime()) % 1000000,
							"account":    "account-" + faker.Username(),
							"repo_owner": "workspace-" + faker.Username(),
							// Missing repo_name
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
				assert.Contains(t, content.Text, "Missing or invalid repo_name parameter")
			})

			t.Run("should handle service error in UpdatePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := int(faker.RandomUnixTime()) % 1000000
				accountName := "account-" + faker.Username()
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				newTitle := "Updated-PR-" + faker.Sentence()
				newDescription := faker.Paragraph()

				// Create expected error
				expectedError := errors.New("failed to update pull request: " + faker.Sentence())

				// Setup mock to return an error
				mockService.EXPECT().
					UpdatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketUpdatePRParams) bool {
						return params.PullRequestID == prID &&
							params.AccountName == accountName &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName &&
							params.Title == newTitle &&
							params.Description == newDescription
					})).
					Return(nil, expectedError)

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
				require.Error(t, err)
				assert.Contains(t, err.Error(), expectedError.Error())
				assert.Nil(t, result)
			})

			t.Run("should handle missing attributes for updating a PR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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
				assert.Contains(t, content.Text, "Missing attributes to update a PR")
			})

			t.Run("should allow just a draft to be updated", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				prID := int(faker.RandomUnixTime()) % 1000000 // Generate a random PR ID within a reasonable range
				draft := rand.IntN(2) == 1

				// Create expected parameters and response
				expectedParams := app.BitbucketUpdatePRParams{
					PullRequestID: prID,
					AccountName:   "account-" + faker.Username(),
					RepoOwner:     "workspace-" + faker.Username(),
					RepoName:      "repo-" + faker.Word(),
					Draft:         lo.ToPtr(draft),
				}

				// Create expected PR response with random values
				expectedPR := &bitbucket.PullRequest{
					ID: prID,
				}

				// Setup mock expectations with any matcher for context
				mockService.EXPECT().
					UpdatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketUpdatePRParams) bool {
						return params.PullRequestID == expectedParams.PullRequestID &&
							*params.Draft == *expectedParams.Draft
					})).
					Return(expectedPR, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"account":    expectedParams.AccountName,
							"repo_owner": expectedParams.RepoOwner,
							"repo_name":  expectedParams.RepoName,
							"draft":      draft,
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
			})

			t.Run("should allow just a title to be updated", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				prID := int(faker.RandomUnixTime()) % 1000000 // Generate a random PR ID within a reasonable range

				// Create expected parameters and response
				expectedParams := app.BitbucketUpdatePRParams{
					PullRequestID: prID,
					AccountName:   "account-" + faker.Username(),
					RepoOwner:     "workspace-" + faker.Username(),
					RepoName:      "repo-" + faker.Word(),
					Title:         "Updated-PR-" + faker.Sentence(),
				}

				// Create expected PR response with random values
				expectedPR := &bitbucket.PullRequest{
					ID: prID,
				}

				// Setup mock expectations with any matcher for context
				mockService.EXPECT().
					UpdatePR(mock.Anything, mock.MatchedBy(func(params app.BitbucketUpdatePRParams) bool {
						return assert.Equal(t, expectedParams.PullRequestID, params.PullRequestID) &&
							assert.Nil(t, params.Draft) &&
							assert.Equal(t, expectedParams.Title, params.Title)
					})).
					Return(expectedPR, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"account":    expectedParams.AccountName,
							"repo_owner": expectedParams.RepoOwner,
							"repo_name":  expectedParams.RepoName,
							"title":      expectedParams.Title,
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
			})
		})

		t.Run("bitbucket_approve_pr", func(t *testing.T) {
			t.Run("should handle ApprovePR call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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
				deps := makeMockDeps(t)
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

			t.Run("should handle missing repo_owner parameter in ApprovePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_approve_pr",
						Arguments: map[string]interface{}{
							"pr_id":   int(faker.RandomUnixTime()) % 1000000,
							"account": "account-" + faker.Username(),
							// Missing repo_owner
							"repo_name": "repo-" + faker.Word(),
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
				assert.Contains(t, content.Text, "Missing or invalid repo_owner parameter")
			})

			t.Run("should handle missing repo_name parameter in ApprovePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_approve_pr",
						Arguments: map[string]interface{}{
							"pr_id":      int(faker.RandomUnixTime()) % 1000000,
							"account":    "account-" + faker.Username(),
							"repo_owner": "workspace-" + faker.Username(),
							// Missing repo_name
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
				assert.Contains(t, content.Text, "Missing or invalid repo_name parameter")
			})

			t.Run("should handle service error in ApprovePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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
				require.Error(t, err)
				assert.Contains(t, err.Error(), expectedError.Error())
				assert.Nil(t, result)
			})
		})

		t.Run("bitbucket_merge_pr", func(t *testing.T) {
			t.Run("should handle MergePR call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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
				deps := makeMockDeps(t)
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

			t.Run("should handle missing repo_name parameter in MergePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_merge_pr",
						Arguments: map[string]interface{}{
							"pr_id":      int(faker.RandomUnixTime()) % 1000000,
							"account":    "account-" + faker.Username(),
							"repo_owner": "workspace-" + faker.Username(),
							// Missing repo_name
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
				assert.Contains(t, content.Text, "Missing or invalid repo_name parameter")
			})

			t.Run("should handle service error in MergePR", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
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
				require.Error(t, err)
				assert.Contains(t, err.Error(), expectedError.Error())
				assert.Nil(t, result)
			})

			t.Run("should handle MergePR with closeSourceBranch set to false", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := int(faker.RandomUnixTime()) % 1000000
				accountName := "account-" + faker.Username()
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				mergeStrategy := "squash"
				commitMessage := "Merge PR: " + faker.Sentence()
				closeSourceBranchStr := "false" // This should set closeSourceBranch to false

				// Create expected parameters
				expectedParams := app.BitbucketMergePRParams{
					PullRequestID:     prID,
					AccountName:       accountName,
					RepoOwner:         repoOwner,
					RepoName:          repoName,
					MergeStrategy:     mergeStrategy,
					Message:           commitMessage,
					CloseSourceBranch: false, // Should be false when closeSourceBranchStr is "false"
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
							"close_source_branch": closeSourceBranchStr,
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
				assert.NotContains(t, content.Text, "source branch was closed") // Should not contain this text
			})
		})

		t.Run("bitbucket_list_pr_tasks", func(t *testing.T) {
			t.Run("should handle ListPRTasks call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				prID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				accountName := "account-" + faker.Username()

				// Create expected parameters
				expectedParams := app.BitbucketListTasksParams{
					PullRequestID: prID,
					RepoOwner:     repoOwner,
					RepoName:      repoName,
					AccountName:   accountName,
				}

				// Create expected tasks response
				createdOn := time.Now()
				updatedOn := time.Now()
				taskID := faker.RandomUnixTime() % 1000000

				expectedTasks := &bitbucket.PaginatedTasks{
					Size:    2,
					Page:    1,
					PageLen: 10,
					Values: []bitbucket.PullRequestCommentTask{
						{
							PullRequestTask: bitbucket.PullRequestTask{
								Task: bitbucket.Task{
									ID:        taskID,
									CreatedOn: createdOn,
									UpdatedOn: updatedOn,
									State:     "RESOLVED",
									Content: &bitbucket.TaskContent{
										Raw:    "Task 1: " + faker.Sentence(),
										Markup: "markdown",
										HTML:   "<p>Task 1: " + faker.Sentence() + "</p>",
									},
									Creator: &bitbucket.Account{
										DisplayName: "User-" + faker.Name(),
										UUID:        faker.UUIDHyphenated(),
									},
									ResolvedOn: time.Now(),
									ResolvedBy: &bitbucket.Account{
										DisplayName: "Resolver-" + faker.Name(),
										UUID:        faker.UUIDHyphenated(),
									},
								},
							},
						},
						{
							PullRequestTask: bitbucket.PullRequestTask{
								Task: bitbucket.Task{
									ID:        taskID + 1,
									CreatedOn: createdOn,
									UpdatedOn: updatedOn,
									State:     "UNRESOLVED",
									Content: &bitbucket.TaskContent{
										Raw:    "Task 2: " + faker.Sentence(),
										Markup: "markdown",
										HTML:   "<p>Task 2: " + faker.Sentence() + "</p>",
									},
									Creator: &bitbucket.Account{
										DisplayName: "User-" + faker.Name(),
										UUID:        faker.UUIDHyphenated(),
									},
								},
							},
						},
					},
				}

				// Setup mock expectations
				mockService.EXPECT().
					ListTasks(mock.Anything, mock.MatchedBy(func(params app.BitbucketListTasksParams) bool {
						return params.PullRequestID == expectedParams.PullRequestID &&
							params.AccountName == expectedParams.AccountName &&
							params.RepoOwner == expectedParams.RepoOwner &&
							params.RepoName == expectedParams.RepoName
					})).
					Return(expectedTasks, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"account":    accountName,
						},
					},
				}

				// Get the handler
				serverTool := controller.newListPRTasksServerTool()
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
				assert.Contains(t, content.Text, fmt.Sprintf("Found %d tasks", expectedTasks.Size))
				assert.Contains(t, content.Text, "RESOLVED")
				assert.Contains(t, content.Text, "UNRESOLVED")
			})

			t.Run("should handle missing pr_id parameter in ListPRTasks", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							// Missing pr_id
							"repo_owner": "workspace-" + faker.Username(),
							"repo_name":  "repo-" + faker.Word(),
							"account":    "account-" + faker.Username(),
						},
					},
				}

				// Get the handler
				serverTool := controller.newListPRTasksServerTool()
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

			t.Run("should handle missing repo_owner parameter in ListPRTasks", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							"pr_id":     int(faker.RandomUnixTime()) % 1000000,
							"repo_name": "repo-" + faker.Word(),
							"account":   "account-" + faker.Username(),
							// Missing repo_owner
						},
					},
				}

				// Get the handler
				serverTool := controller.newListPRTasksServerTool()
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

			t.Run("should handle missing repo_name parameter in ListPRTasks", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							"pr_id":      int(faker.RandomUnixTime()) % 1000000,
							"repo_owner": "workspace-" + faker.Username(),
							"account":    "account-" + faker.Username(),
							// Missing repo_name
						},
					},
				}

				// Get the handler
				serverTool := controller.newListPRTasksServerTool()
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
				assert.Contains(t, content.Text, "Missing or invalid repo_name parameter")
			})

			t.Run("should handle service error in ListPRTasks", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := int(faker.RandomUnixTime()) % 1000000
				accountName := "account-" + faker.Username()
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()

				// Create expected error
				expectedError := errors.New("failed to list tasks: " + faker.Sentence())

				// Setup mock to return an error
				mockService.EXPECT().
					ListTasks(mock.Anything, mock.MatchedBy(func(params app.BitbucketListTasksParams) bool {
						return params.PullRequestID == prID &&
							params.AccountName == accountName &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName
					})).
					Return(nil, expectedError)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"account":    accountName,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
						},
					},
				}

				// Get the handler
				serverTool := controller.newListPRTasksServerTool()
				handler := serverTool.Handler

				// Act
				result, err := handler(ctx, request)

				// Assert
				require.Error(t, err)
				assert.Contains(t, err.Error(), expectedError.Error())
				assert.Nil(t, result)
			})
		})

		t.Run("bitbucket_update_pr_task", func(t *testing.T) {
			t.Run("should handle UpdatePRTask call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				prID := int(faker.RandomUnixTime()) % 1000000
				taskID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				accountName := "account-" + faker.Username()
				newContent := "Updated task: " + faker.Sentence()
				newState := "RESOLVED"

				// Create expected parameters
				expectedParams := app.BitbucketUpdateTaskParams{
					PullRequestID: prID,
					TaskID:        taskID,
					RepoOwner:     repoOwner,
					RepoName:      repoName,
					AccountName:   accountName,
					Content:       newContent,
					State:         newState,
				}

				// Create expected task response
				createdOn := time.Now()
				updatedOn := time.Now()

				expectedTask := &bitbucket.PullRequestCommentTask{
					PullRequestTask: bitbucket.PullRequestTask{
						Task: bitbucket.Task{
							ID:        int64(taskID),
							CreatedOn: createdOn,
							UpdatedOn: updatedOn,
							State:     newState,
							Content: &bitbucket.TaskContent{
								Raw:    newContent,
								Markup: "markdown",
								HTML:   "<p>" + newContent + "</p>",
							},
							Creator: &bitbucket.Account{
								DisplayName: "User-" + faker.Name(),
								UUID:        faker.UUIDHyphenated(),
							},
							ResolvedOn: time.Now(),
							ResolvedBy: &bitbucket.Account{
								DisplayName: "Resolver-" + faker.Name(),
								UUID:        faker.UUIDHyphenated(),
							},
						},
					},
				}

				// Setup mock expectations
				mockService.EXPECT().
					UpdateTask(mock.Anything, mock.MatchedBy(func(params app.BitbucketUpdateTaskParams) bool {
						return params.PullRequestID == expectedParams.PullRequestID &&
							params.TaskID == expectedParams.TaskID &&
							params.AccountName == expectedParams.AccountName &&
							params.RepoOwner == expectedParams.RepoOwner &&
							params.RepoName == expectedParams.RepoName &&
							params.Content == expectedParams.Content &&
							params.State == expectedParams.State
					})).
					Return(expectedTask, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr_task",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"task_id":    taskID,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"account":    accountName,
							"content":    newContent,
							"state":      newState,
						},
					},
				}

				// Get the handler
				serverTool := controller.newUpdatePRTaskServerTool()
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
				assert.Contains(t, content.Text, fmt.Sprintf("Updated task #%d", taskID))
				assert.Contains(t, content.Text, "RESOLVED")
			})

			t.Run("should handle update with only state change", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data
				prID := int(faker.RandomUnixTime()) % 1000000
				taskID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				accountName := "account-" + faker.Username()
				newState := "RESOLVED"

				// Create expected parameters (no content, only state)
				expectedParams := app.BitbucketUpdateTaskParams{
					PullRequestID: prID,
					TaskID:        taskID,
					RepoOwner:     repoOwner,
					RepoName:      repoName,
					AccountName:   accountName,
					State:         newState,
				}

				// Create expected task response
				createdOn := time.Now()
				updatedOn := time.Now()
				existingContent := "Task content: " + faker.Sentence()

				expectedTask := &bitbucket.PullRequestCommentTask{
					PullRequestTask: bitbucket.PullRequestTask{
						Task: bitbucket.Task{
							ID:        int64(taskID),
							CreatedOn: createdOn,
							UpdatedOn: updatedOn,
							State:     newState,
							Content: &bitbucket.TaskContent{
								Raw:    existingContent,
								Markup: "markdown",
								HTML:   "<p>" + existingContent + "</p>",
							},
							Creator: &bitbucket.Account{
								DisplayName: "User-" + faker.Name(),
								UUID:        faker.UUIDHyphenated(),
							},
							ResolvedOn: time.Now(),
							ResolvedBy: &bitbucket.Account{
								DisplayName: "Resolver-" + faker.Name(),
								UUID:        faker.UUIDHyphenated(),
							},
						},
					},
				}

				// Setup mock expectations
				mockService.EXPECT().
					UpdateTask(mock.Anything, mock.MatchedBy(func(params app.BitbucketUpdateTaskParams) bool {
						return params.PullRequestID == expectedParams.PullRequestID &&
							params.TaskID == expectedParams.TaskID &&
							params.AccountName == expectedParams.AccountName &&
							params.RepoOwner == expectedParams.RepoOwner &&
							params.RepoName == expectedParams.RepoName &&
							params.Content == "" &&
							params.State == expectedParams.State
					})).
					Return(expectedTask, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr_task",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"task_id":    taskID,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"account":    accountName,
							"state":      newState,
						},
					},
				}

				// Get the handler
				serverTool := controller.newUpdatePRTaskServerTool()
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
				assert.Contains(t, content.Text, fmt.Sprintf("Updated task #%d", taskID))
				assert.Contains(t, content.Text, "RESOLVED")
			})

			t.Run("should handle service error in UpdatePRTask", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data
				prID := int(faker.RandomUnixTime()) % 1000000
				taskID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				newContent := "Updated content"

				// Setup mock to return an error
				expectedError := errors.New("task update failed")
				mockService.EXPECT().
					UpdateTask(mock.Anything, mock.Anything).
					Return(nil, expectedError)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr_task",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"task_id":    taskID,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"content":    newContent,
						},
					},
				}

				// Get the handler
				serverTool := controller.newUpdatePRTaskServerTool()
				handler := serverTool.Handler

				// Act
				_, err := handler(ctx, request)

				// Assert
				require.Error(t, err)
				assert.Contains(t, err.Error(), "task update failed")
			})

			t.Run("should handle missing pr_id parameter in UpdatePRTask", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr_task",
						Arguments: map[string]interface{}{
							// Missing pr_id
							"task_id":    int(faker.RandomUnixTime()) % 1000000,
							"repo_owner": "workspace-" + faker.Username(),
							"repo_name":  "repo-" + faker.Word(),
							"content":    "New content",
						},
					},
				}

				// Get the handler
				serverTool := controller.newUpdatePRTaskServerTool()
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

			t.Run("should handle missing task_id parameter in UpdatePRTask", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr_task",
						Arguments: map[string]interface{}{
							"pr_id": int(faker.RandomUnixTime()) % 1000000,
							// Missing task_id
							"repo_owner": "workspace-" + faker.Username(),
							"repo_name":  "repo-" + faker.Word(),
							"content":    "New content",
						},
					},
				}

				// Get the handler
				serverTool := controller.newUpdatePRTaskServerTool()
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
				assert.Contains(t, content.Text, "Missing or invalid task_id parameter")
			})

			t.Run("should handle missing both content and state parameters in UpdatePRTask", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_update_pr_task",
						Arguments: map[string]interface{}{
							"pr_id":      int(faker.RandomUnixTime()) % 1000000,
							"task_id":    int(faker.RandomUnixTime()) % 1000000,
							"repo_owner": "workspace-" + faker.Username(),
							"repo_name":  "repo-" + faker.Word(),
							// Missing both content and state
						},
					},
				}

				// Get the handler
				serverTool := controller.newUpdatePRTaskServerTool()
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
				assert.Contains(t, content.Text, "Either content or state must be provided")
			})

			t.Run("should display actual task IDs in ListPRTasks output", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values using faker
				prID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				accountName := "account-" + faker.Username()

				// Create tasks with specific IDs for testing
				taskID1 := 100 + rand.Int64N(10000)
				taskID2 := 100 + rand.Int64N(10000)
				task1Content := "Task: " + faker.Sentence()
				task2Content := "Task: " + faker.Sentence()
				creator1Name := "User-" + faker.Name()
				creator2Name := "User-" + faker.Name()

				expectedTasks := &bitbucket.PaginatedTasks{
					Size:    2,
					Page:    1,
					PageLen: 10,
					Values: []bitbucket.PullRequestCommentTask{
						{
							PullRequestTask: bitbucket.PullRequestTask{
								Task: bitbucket.Task{
									ID:    taskID1,
									State: "RESOLVED",
									Content: &bitbucket.TaskContent{
										Raw: task1Content,
									},
									Creator: &bitbucket.Account{
										DisplayName: creator1Name,
									},
								},
							},
						},
						{
							PullRequestTask: bitbucket.PullRequestTask{
								Task: bitbucket.Task{
									ID:    taskID2,
									State: "UNRESOLVED",
									Content: &bitbucket.TaskContent{
										Raw: task2Content,
									},
									Creator: &bitbucket.Account{
										DisplayName: creator2Name,
									},
								},
							},
						},
					},
				}

				// Setup mock expectations
				mockService.EXPECT().
					ListTasks(ctx, mock.MatchedBy(func(params app.BitbucketListTasksParams) bool {
						return params.PullRequestID == prID &&
							params.AccountName == accountName &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName
					})).
					Return(expectedTasks, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"account":    accountName,
						},
					},
				}

				// Act
				serverTool := controller.newListPRTasksServerTool()
				result, err := serverTool.Handler(ctx, request)

				// Assert
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.False(t, result.IsError)

				// Verify the content of the result
				content, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok, "Result content should be text content")

				// Verify task count is shown correctly
				assert.Contains(t, content.Text, fmt.Sprintf("Found %d tasks", expectedTasks.Size))

				// Verify that the task IDs are displayed, not indices
				assert.Contains(t, content.Text, fmt.Sprintf("Task #%d", taskID1))
				assert.Contains(t, content.Text, fmt.Sprintf("Task #%d", taskID2))

				// Verify task content is shown correctly
				assert.Contains(t, content.Text, task1Content)
				assert.Contains(t, content.Text, task2Content)

				// Verify that indices are not used in the output
				assert.NotContains(t, content.Text, "1. [")
				assert.NotContains(t, content.Text, "2. [")
			})

			t.Run("should handle nil Creator in ListPRTasks", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := int(faker.RandomUnixTime()) % 1000000
				repoOwner := "workspace-" + faker.Username()
				repoName := "repo-" + faker.Word()
				accountName := "account-" + faker.Username()

				// Create tasks with specific IDs for testing
				taskID1 := 100 + rand.Int64N(10000)
				taskID2 := 100 + rand.Int64N(10000)
				task1Content := "Task: " + faker.Sentence()
				task2Content := "Task: " + faker.Sentence()

				expectedTasks := &bitbucket.PaginatedTasks{
					Size:    2,
					Page:    1,
					PageLen: 10,
					Values: []bitbucket.PullRequestCommentTask{
						{
							PullRequestTask: bitbucket.PullRequestTask{
								Task: bitbucket.Task{
									ID:    taskID1,
									State: "RESOLVED",
									Content: &bitbucket.TaskContent{
										Raw: task1Content,
									},
									Creator: nil, // Explicitly set Creator to nil for testing
								},
							},
						},
						{
							PullRequestTask: bitbucket.PullRequestTask{
								Task: bitbucket.Task{
									ID:    taskID2,
									State: "UNRESOLVED",
									Content: &bitbucket.TaskContent{
										Raw: task2Content,
									},
									// Creator is nil by default
								},
							},
						},
					},
				}

				// Setup mock expectations
				mockService.EXPECT().
					ListTasks(ctx, mock.MatchedBy(func(params app.BitbucketListTasksParams) bool {
						return params.PullRequestID == prID &&
							params.AccountName == accountName &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName
					})).
					Return(expectedTasks, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_list_pr_tasks",
						Arguments: map[string]interface{}{
							"pr_id":      prID,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"account":    accountName,
						},
					},
				}

				// Act
				serverTool := controller.newListPRTasksServerTool()
				result, err := serverTool.Handler(ctx, request)

				// Assert
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.False(t, result.IsError)

				// Verify the content of the result
				content, ok := result.Content[0].(mcp.TextContent)
				require.True(t, ok, "Result content should be text content")

				// Verify task count is shown correctly
				assert.Contains(t, content.Text, fmt.Sprintf("Found %d tasks", expectedTasks.Size))

				// Verify that the task IDs are displayed
				assert.Contains(t, content.Text, fmt.Sprintf("Task #%d", taskID1))
				assert.Contains(t, content.Text, fmt.Sprintf("Task #%d", taskID2))

				// Verify task content is shown correctly
				assert.Contains(t, content.Text, task1Content)
				assert.Contains(t, content.Text, task2Content)

				// Verify that the response doesn't crash due to nil Creator
				assert.Contains(t, content.Text, "unknown user")
			})
		})

		t.Run("bitbucket_create_pr_task", func(t *testing.T) {
			t.Run("should handle CreatePRTask call successfully", func(t *testing.T) {
				// Arrange
				deps := makeMockDeps(t)
				mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
				controller := NewBitbucketController(deps)
				ctx := t.Context()

				// Create test data with randomized values
				prID := 123
				taskContent := "Task: Review code formatting"
				repoOwner := "workspace-abc"
				repoName := "repo-xyz"
				account := "account-1"
				state := "UNRESOLVED"

				// Create expected task response with proper structure
				expectedTask := &bitbucket.PullRequestCommentTask{
					PullRequestTask: bitbucket.PullRequestTask{
						Task: bitbucket.Task{
							ID:    789,
							State: state,
							Content: &bitbucket.TaskContent{
								Raw: taskContent,
							},
							Creator: &bitbucket.Account{
								DisplayName: "Test User",
							},
						},
					},
				}

				// Setup mock expectations with looser matching criteria
				mockService.EXPECT().
					CreateTask(mock.Anything, mock.MatchedBy(func(params app.BitbucketCreateTaskParams) bool {
						return params.PullRequestID == prID &&
							params.Content == taskContent &&
							params.RepoOwner == repoOwner &&
							params.RepoName == repoName &&
							params.AccountName == account &&
							params.State == state
					})).
					Return(expectedTask, nil)

				// Create the request
				request := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name: "bitbucket_create_pr_task",
						Arguments: map[string]interface{}{
							"pr_id":      float64(prID),
							"content":    taskContent,
							"repo_owner": repoOwner,
							"repo_name":  repoName,
							"account":    account,
							"state":      state,
						},
					},
				}

				// Get the handler
				serverTool := controller.newCreatePRTaskServerTool()
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
				assert.Contains(t, content.Text, "Created task")
				assert.Contains(t, content.Text, taskContent)
				assert.Contains(t, content.Text, fmt.Sprintf("on PR #%d", prID))
			})
		})
	})

	t.Run("bitbucket_get_pr_diffstat", func(t *testing.T) {
		t.Run("happy path: returns diffstat and summary", func(t *testing.T) {
			deps := makeMockDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			// Arrange: randomized test data
			prID := int(faker.RandomUnixTime()) % 1000000
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			accountName := "account-" + faker.Username()
			file1 := faker.Word() + ".go"
			file2 := faker.Word() + ".go"

			expectedDiffstat := &app.PaginatedDiffStat{
				Size: 2,
				Values: []bitbucket.DiffStat{
					{
						Status:       "modified",
						LinesAdded:   10,
						LinesRemoved: 2,
						Old:          file1,
						New:          file1,
					},
					{
						Status:       "added",
						LinesAdded:   42,
						LinesRemoved: 0,
						Old:          "",
						New:          file2,
					},
				},
			}

			mockService.EXPECT().
				GetPRDiffStat(ctx, mock.MatchedBy(func(params app.BitbucketGetPRDiffStatParams) bool {
					return params.PullRequestID == prID &&
						params.RepoOwner == repoOwner &&
						params.RepoName == repoName &&
						params.AccountName == accountName
				})).
				Return(expectedDiffstat, nil)

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_get_pr_diffstat",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
						"account":    accountName,
					},
				},
			}
			serverTool := controller.newGetPRDiffstatServerTool()
			handler := serverTool.Handler

			// Act
			result, err := handler(ctx, request)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.False(t, result.IsError)
			require.Len(t, result.Content, 2)

			summary, ok := result.Content[0].(mcp.TextContent)
			require.True(t, ok)
			assert.Contains(t, summary.Text, "Diffstat for PR #")
			assert.Contains(t, summary.Text, "2 files changed")

			jsonContent, ok := result.Content[1].(mcp.TextContent)
			require.True(t, ok)
			var parsed app.PaginatedDiffStat
			err = json.Unmarshal([]byte(jsonContent.Text), &parsed)
			require.NoError(t, err)
			assert.Equal(t, *expectedDiffstat, parsed)
		})

		t.Run("missing required parameter: pr_id", func(t *testing.T) {
			deps := makeMockDeps(t)
			controller := NewBitbucketController(deps)
			ctx := t.Context()
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			accountName := "account-" + faker.Username()

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_get_pr_diffstat",
					Arguments: map[string]interface{}{
						// "pr_id" is missing
						"repo_owner": repoOwner,
						"repo_name":  repoName,
						"account":    accountName,
					},
				},
			}
			serverTool := controller.newGetPRDiffstatServerTool()
			handler := serverTool.Handler

			result, err := handler(ctx, request)

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.True(t, result.IsError)
			require.Len(t, result.Content, 1)
			content, ok := result.Content[0].(mcp.TextContent)
			require.True(t, ok)
			assert.Contains(t, content.Text, "Missing or invalid pr_id")
		})

		t.Run("service error is returned", func(t *testing.T) {
			deps := makeMockDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			prID := int(faker.RandomUnixTime()) % 1000000
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			accountName := "account-" + faker.Username()
			expectedErr := errors.New("bitbucket service failure: " + faker.Sentence())

			mockService.EXPECT().
				GetPRDiffStat(ctx, mock.Anything).
				Return(nil, expectedErr)

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_get_pr_diffstat",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
						"account":    accountName,
					},
				},
			}
			serverTool := controller.newGetPRDiffstatServerTool()
			handler := serverTool.Handler

			result, err := handler(ctx, request)

			require.Error(t, err)
			assert.Contains(t, err.Error(), expectedErr.Error())
			assert.Nil(t, result)
		})

		t.Run("returns empty diffstat (zero files changed)", func(t *testing.T) {
			deps := makeMockDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			prID := int(faker.RandomUnixTime()) % 1000000
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			accountName := "account-" + faker.Username()

			emptyDiffstat := &app.PaginatedDiffStat{
				Size:   0,
				Values: []bitbucket.DiffStat{},
			}

			mockService.EXPECT().
				GetPRDiffStat(ctx, mock.Anything).
				Return(emptyDiffstat, nil)

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_get_pr_diffstat",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
						"account":    accountName,
					},
				},
			}
			serverTool := controller.newGetPRDiffstatServerTool()
			handler := serverTool.Handler

			result, err := handler(ctx, request)

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.False(t, result.IsError)
			require.Len(t, result.Content, 2)

			summary, ok := result.Content[0].(mcp.TextContent)
			require.True(t, ok)
			assert.Contains(t, summary.Text, "Diffstat for PR #")
			assert.Contains(t, summary.Text, "0 files changed")

			jsonContent, ok := result.Content[1].(mcp.TextContent)
			require.True(t, ok)
			var parsed app.PaginatedDiffStat
			err = json.Unmarshal([]byte(jsonContent.Text), &parsed)
			require.NoError(t, err)
			assert.Equal(t, *emptyDiffstat, parsed)
		})

		t.Run("handles optional file_paths parameter", func(t *testing.T) {
			deps := makeMockDeps(t)
			mockService := mocks.GetMock[*MockbitbucketService](t, deps.BitbucketService)
			controller := NewBitbucketController(deps)
			ctx := t.Context()

			prID := int(faker.RandomUnixTime()) % 1000000
			repoOwner := "workspace-" + faker.Username()
			repoName := "repo-" + faker.Word()
			accountName := "account-" + faker.Username()
			filePaths := []string{faker.Word() + ".go", faker.Word() + ".py"}

			expectedDiffstat := &app.PaginatedDiffStat{
				Size: 2,
				Values: []bitbucket.DiffStat{
					{
						Status:       "modified",
						LinesAdded:   10,
						LinesRemoved: 2,
						Old:          filePaths[0],
						New:          filePaths[0],
					},
					{
						Status:       "added",
						LinesAdded:   42,
						LinesRemoved: 0,
						Old:          "",
						New:          filePaths[1],
					},
				},
			}

			mockService.EXPECT().
				GetPRDiffStat(ctx, mock.MatchedBy(func(params app.BitbucketGetPRDiffStatParams) bool {
					return params.PullRequestID == prID &&
						params.RepoOwner == repoOwner &&
						params.RepoName == repoName &&
						params.AccountName == accountName
				})).
				Return(expectedDiffstat, nil)

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bitbucket_get_pr_diffstat",
					Arguments: map[string]interface{}{
						"pr_id":      prID,
						"repo_owner": repoOwner,
						"repo_name":  repoName,
						"account":    accountName,
						"file_paths": filePaths,
					},
				},
			}
			serverTool := controller.newGetPRDiffstatServerTool()
			handler := serverTool.Handler

			result, err := handler(ctx, request)

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.False(t, result.IsError)
			require.Len(t, result.Content, 2)

			summary, ok := result.Content[0].(mcp.TextContent)
			require.True(t, ok)
			assert.Contains(t, summary.Text, "Diffstat for PR #")
			assert.Contains(t, summary.Text, "2 files changed")

			jsonContent, ok := result.Content[1].(mcp.TextContent)
			require.True(t, ok)
			var parsed app.PaginatedDiffStat
			err = json.Unmarshal([]byte(jsonContent.Text), &parsed)
			require.NoError(t, err)
			assert.Equal(t, *expectedDiffstat, parsed)
		})
	})
}
