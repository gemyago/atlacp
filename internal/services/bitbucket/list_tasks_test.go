package bitbucket

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/atlacp/internal/diag"
	httpservices "github.com/gemyago/atlacp/internal/services/http"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListPullRequestTasks(t *testing.T) {
	// Explicit t parameter for makeMockDeps
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

	mockTokenProvider := &MockTokenProvider{}
	tokenError := errors.New("failed to get token")

	t.Run("success with all parameters", func(t *testing.T) {
		// Arrange - Use randomized data
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 123
		queryParam := "test-" + faker.Word()
		sortParam := "created_on"
		pageLen := 50

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks", workspace, repoSlug, pullReqID), r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, queryParam, r.URL.Query().Get("q"))
			assert.Equal(t, sortParam, r.URL.Query().Get("sort"))
			assert.Equal(t, "50", r.URL.Query().Get("pagelen"))

			// Return successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"pagelen": 50,
				"size": 2,
				"values": [
					{
						"id": 1001,
						"created_on": "2023-01-01T12:00:00.000Z",
						"updated_on": "2023-01-01T13:00:00.000Z",
						"state": "RESOLVED",
						"content": {
							"raw": "Task 1 description",
							"markup": "markdown",
							"html": "<p>Task 1 description</p>"
						},
						"creator": {
							"display_name": "John Doe",
							"uuid": "abc123"
						},
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/repositories/workspace/repo/pullrequests/123/tasks/1001"
							}
						},
						"comment": {
							"id": 5001
						}
					},
					{
						"id": 1002,
						"created_on": "2023-01-02T12:00:00.000Z",
						"updated_on": "2023-01-02T13:00:00.000Z",
						"state": "UNRESOLVED",
						"content": {
							"raw": "Task 2 description",
							"markup": "markdown",
							"html": "<p>Task 2 description</p>"
						},
						"creator": {
							"display_name": "Jane Smith",
							"uuid": "def456"
						},
						"links": {
							"self": {
								"href": "https://api.bitbucket.org/2.0/repositories/workspace/repo/pullrequests/123/tasks/1002"
							}
						},
						"comment": {
							"id": 5002
						}
					}
				],
				"page": 1,
				"next": "https://api.bitbucket.org/2.0/repositories/workspace/repo/pullrequests/123/tasks?page=2"
			}`)
		}))
		defer server.Close()

		deps := makeMockDeps(t, server.URL)
		client := NewClient(deps)

		mockTokenProvider.Token = "test-token"
		params := ListPullRequestTasksParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Query:     queryParam,
			Sort:      sortParam,
			PageLen:   pageLen,
		}

		// Act
		result, err := client.ListPullRequestTasks(t.Context(), mockTokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 50, result.PageLen)
		assert.Equal(t, 2, result.Size)
		assert.Len(t, result.Values, 2)

		// Verify first task
		task1 := result.Values[0]
		assert.Equal(t, int64(1001), task1.ID)
		assert.Equal(t, "RESOLVED", task1.State)
		assert.Equal(t, "Task 1 description", task1.Content.Raw)
		assert.Equal(t, "John Doe", task1.Creator.DisplayName)
		assert.Equal(t, int64(5001), task1.Comment.ID)

		// Verify second task
		task2 := result.Values[1]
		assert.Equal(t, int64(1002), task2.ID)
		assert.Equal(t, "UNRESOLVED", task2.State)
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Arrange - Use randomized data
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 123

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks", workspace, repoSlug, pullReqID), r.URL.Path)
			assert.Empty(t, r.URL.Query().Get("q"))
			assert.Empty(t, r.URL.Query().Get("sort"))
			assert.Empty(t, r.URL.Query().Get("pagelen"))

			// Return successful response with minimal data
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"values": [
					{
						"id": 1001,
						"created_on": "2023-01-01T12:00:00.000Z",
						"updated_on": "2023-01-01T13:00:00.000Z",
						"state": "UNRESOLVED",
						"content": {
							"raw": "Task description"
						},
						"creator": {
							"display_name": "John Doe"
						}
					}
				]
			}`)
		}))
		defer server.Close()

		deps := makeMockDeps(t, server.URL)
		client := NewClient(deps)

		mockTokenProvider.Token = "test-token"
		params := ListPullRequestTasksParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
		}

		// Act
		result, err := client.ListPullRequestTasks(t.Context(), mockTokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Values, 1)

		task := result.Values[0]
		assert.Equal(t, int64(1001), task.ID)
		assert.Equal(t, "UNRESOLVED", task.State)
		assert.Equal(t, "Task description", task.Content.Raw)
		assert.Equal(t, "John Doe", task.Creator.DisplayName)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 123

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{
				"type": "error",
				"error": {
					"message": "Pull request not found"
				}
			}`)
		}))
		defer server.Close()

		deps := makeMockDeps(t, server.URL)
		client := NewClient(deps)

		mockTokenProvider.Token = "test-token"
		params := ListPullRequestTasksParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
		}

		// Act
		result, err := client.ListPullRequestTasks(t.Context(), mockTokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "list pull request tasks failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 123

		deps := makeMockDeps(t, "https://api.example.com")
		client := NewClient(deps)

		mockTokenProvider.Err = tokenError
		params := ListPullRequestTasksParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
		}

		// Act
		result, err := client.ListPullRequestTasks(t.Context(), mockTokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		expectedError := fmt.Errorf("failed to get token: %w", tokenError)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}
