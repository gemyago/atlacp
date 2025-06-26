package bitbucket

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_UpdateTask(t *testing.T) {
	makeMockDeps := func(t *testing.T, baseURL string) ClientDeps {
		rootLogger := diag.RootTestLogger().With("test", t.Name())
		return ClientDeps{
			ClientFactory: httpservices.NewClientFactory(httpservices.ClientFactoryDeps{
				RootLogger: rootLogger,
			}),
			RootLogger: rootLogger,
			BaseURL:    baseURL,
		}
	}

	t.Run("success with all parameters", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 1000 + rand.IntN(9000)
		taskID := 100 + rand.IntN(900)
		updatedContent := "Updated task: " + faker.Sentence()
		taskState := "RESOLVED"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
				workspace, repoSlug, pullReqID, taskID), r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"created_on": "2023-01-01T12:00:00Z",
				"updated_on": "2023-01-02T12:00:00Z", 
				"state": "RESOLVED",
				"content": {
					"raw": "%s",
					"html": "<p>%s</p>",
					"markup": "markdown"
				},
				"creator": {
					"type": "user",
					"display_name": "Test User",
					"uuid": "{1234567890}",
					"account_id": "123456:7890"
				},
				"resolved_on": "2023-01-02T12:00:00Z",
				"resolved_by": {
					"type": "user",
					"display_name": "Resolver User",
					"uuid": "{0987654321}",
					"account_id": "098765:4321"
				},
				"links": {
					"self": { 
						"href": "https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests/%d/tasks/%d"
					}
				}
			}`, taskID, updatedContent, updatedContent, workspace, repoSlug, pullReqID, taskID)
		}))
		defer mockServer.Close()

		deps := makeMockDeps(t, mockServer.URL)
		client := NewClient(deps)

		// Act
		task, err := client.UpdateTask(t.Context(), mockTokenProvider, UpdateTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
			Content:   updatedContent,
			State:     taskState,
		})

		// Assert
		require.NoError(t, err)
		require.NotNil(t, task)
		assert.Equal(t, int64(taskID), task.ID)
		assert.Equal(t, updatedContent, task.Content.Raw)
		assert.Equal(t, "RESOLVED", task.State)
		assert.NotZero(t, task.CreatedOn)
		assert.NotZero(t, task.UpdatedOn)
		assert.NotZero(t, task.ResolvedOn)
		assert.NotNil(t, task.ResolvedBy)
	})

	t.Run("success with only task state update", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 1000 + rand.IntN(9000)
		taskID := 100 + rand.IntN(900)
		originalContent := "Original task: " + faker.Sentence()
		taskState := "RESOLVED"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"created_on": "2023-01-01T12:00:00Z",
				"updated_on": "2023-01-02T12:00:00Z", 
				"state": "RESOLVED",
				"content": {
					"raw": "%s",
					"html": "<p>%s</p>",
					"markup": "markdown"
				},
				"creator": {
					"type": "user",
					"display_name": "Test User",
					"uuid": "{1234567890}",
					"account_id": "123456:7890"
				},
				"resolved_on": "2023-01-02T12:00:00Z",
				"resolved_by": {
					"type": "user",
					"display_name": "Resolver User",
					"uuid": "{0987654321}",
					"account_id": "098765:4321"
				}
			}`, taskID, originalContent, originalContent)
		}))
		defer mockServer.Close()

		deps := makeMockDeps(t, mockServer.URL)
		client := NewClient(deps)

		// Act - only updating state, not content
		task, err := client.UpdateTask(t.Context(), mockTokenProvider, UpdateTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
			State:     taskState, // Only updating state
		})

		// Assert
		require.NoError(t, err)
		require.NotNil(t, task)
		assert.Equal(t, int64(taskID), task.ID)
		assert.Equal(t, originalContent, task.Content.Raw)
		assert.Equal(t, "RESOLVED", task.State)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 1000 + rand.IntN(9000)
		taskID := 100 + rand.IntN(900)
		taskState := "RESOLVED"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{
				"type": "error",
				"error": {
					"message": "Invalid task state"
				}
			}`)
		}))
		defer mockServer.Close()

		deps := makeMockDeps(t, mockServer.URL)
		client := NewClient(deps)

		// Act
		task, err := client.UpdateTask(t.Context(), mockTokenProvider, UpdateTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
			State:     taskState,
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, task)
		assert.ErrorContains(t, err, "update task failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 1000 + rand.IntN(9000)
		taskID := 100 + rand.IntN(900)
		tokenError := errors.New(faker.Sentence())

		mockTokenProvider := &MockTokenProvider{
			Err: tokenError,
		}

		deps := makeMockDeps(t, "http://example.com")
		client := NewClient(deps)

		// Act
		task, err := client.UpdateTask(t.Context(), mockTokenProvider, UpdateTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, task)
		expectedError := fmt.Errorf("failed to get token: %w", tokenError)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}
