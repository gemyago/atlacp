package bitbucket

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetTask(t *testing.T) {
	t.Run("success with all fields", func(t *testing.T) {
		// Arrange - Use randomized data
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 100 + rand.IntN(10000)
		taskID := 100 + rand.IntN(10000)
		displayName := faker.Name()
		uuid := faker.UUIDHyphenated()
		accountID := strconv.Itoa(10000000 + rand.IntN(90000000))
		nickname := "user-" + faker.Username()
		taskDescription := "Task-" + faker.Sentence()
		commentID := int64(100 + rand.IntN(10000))

		mockTokenProvider := &MockTokenProvider{
			TokenType:  faker.Word(),
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "GET", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
				workspace, repoSlug, pullReqID, taskID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, mockTokenProvider.TokenType+" "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"created_on": "2023-01-01T12:00:00.000Z",
				"updated_on": "2023-01-01T13:00:00.000Z",
				"state": "RESOLVED",
				"content": {
					"raw": "%s",
					"markup": "markdown",
					"html": "<p>%s</p>"
				},
				"creator": {
					"display_name": "%s",
					"uuid": "%s",
					"account_id": "%s",
					"nickname": "%s",
					"links": {
						"avatar": {
							"href": "https://avatar.url/%s"
						}
					}
				},
				"resolved_on": "2023-01-02T12:00:00.000Z",
				"resolved_by": {
					"display_name": "%s",
					"uuid": "%s"
				},
				"links": {
					"self": {
						"href": "https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests/%d/tasks/%d"
					},
					"html": {
						"href": "https://bitbucket.org/%s/%s/pull-requests/%d/tasks/%d"
					}
				},
				"comment": {
					"id": %d,
					"content": {
						"raw": "Comment related to %s"
					},
					"user": {
						"display_name": "%s"
					}
				}
			}`,
				taskID,
				taskDescription,
				taskDescription,
				displayName,
				uuid,
				accountID,
				nickname,
				uuid,
				faker.Name(),
				faker.UUIDHyphenated(),
				workspace,
				repoSlug,
				pullReqID,
				taskID,
				workspace,
				repoSlug,
				pullReqID,
				taskID,
				commentID,
				taskDescription,
				displayName)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := GetTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
		}

		// Act
		result, err := client.GetTask(t.Context(), mockTokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify task details
		assert.Equal(t, int64(taskID), result.ID)
		assert.Equal(t, "RESOLVED", result.State)
		assert.Equal(t, taskDescription, result.Content.Raw)
		assert.Equal(t, "markdown", result.Content.Markup)
		assert.Equal(t, "<p>"+taskDescription+"</p>", result.Content.HTML)
		assert.Equal(t, displayName, result.Creator.DisplayName)
		assert.Equal(t, uuid, result.Creator.UUID)
		assert.Equal(t, accountID, result.Creator.AccountID)
		assert.Equal(t, nickname, result.Creator.Nickname)
		assert.Contains(t, result.Creator.Links.Avatar.Href, uuid)

		// Verify links
		assert.Contains(t, result.Links.Self.Href, fmt.Sprintf("pullrequests/%d/tasks/%d", pullReqID, taskID))
		assert.Contains(t, result.Links.HTML.Href, fmt.Sprintf("pull-requests/%d/tasks/%d", pullReqID, taskID))

		// Verify comment
		assert.Equal(t, commentID, result.Comment.ID)
		assert.Contains(t, result.Comment.Content.Raw, taskDescription)
		assert.Equal(t, displayName, result.Comment.User.DisplayName)
	})

	t.Run("success with required fields only", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 100 + rand.IntN(10000)
		taskID := 100 + rand.IntN(10000)
		displayName := faker.Name()
		taskDescription := "Task-" + faker.Sentence()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  faker.Word(),
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "GET", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
				workspace, repoSlug, pullReqID, taskID)
			assert.Equal(t, expectedPath, r.URL.Path)

			// Return minimal success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"created_on": "2023-01-01T12:00:00.000Z",
				"updated_on": "2023-01-01T13:00:00.000Z",
				"state": "UNRESOLVED",
				"content": {
					"raw": "%s"
				},
				"creator": {
					"display_name": "%s"
				}
			}`, taskID, taskDescription, displayName)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := GetTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
		}

		// Act
		result, err := client.GetTask(t.Context(), mockTokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, int64(taskID), result.ID)
		assert.Equal(t, "UNRESOLVED", result.State)
		assert.Equal(t, taskDescription, result.Content.Raw)
		assert.Equal(t, displayName, result.Creator.DisplayName)

		// Optional fields should be empty or nil
		assert.Empty(t, result.Content.Markup)
		assert.Empty(t, result.Content.HTML)
		assert.Empty(t, result.Creator.UUID)
		assert.Nil(t, result.Comment)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Arrange
		workspace := "workspace-" + faker.Word()
		repoSlug := "repo-" + faker.Word()
		pullReqID := 100 + rand.IntN(10000)
		taskID := 100 + rand.IntN(10000)
		errorMessage := "Task-" + faker.Word() + " not found"

		mockTokenProvider := &MockTokenProvider{
			TokenType:  faker.Word(),
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a 404 Not Found error
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{
				"type": "error",
				"error": {
					"message": "%s"
				}
			}`, errorMessage)
		}))
		defer server.Close()

		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		params := GetTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			TaskID:    taskID,
		}

		// Act
		result, err := client.GetTask(t.Context(), mockTokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get task failed")
	})

	t.Run("handles token error", func(t *testing.T) {
		// Arrange
		tokenErr := errors.New("token error")
		mockTokenProvider := &MockTokenProvider{
			Err: tokenErr,
		}

		deps := makeMockDepsWithTestName(t, "http://dummy-url")
		client := NewClient(deps)

		params := GetTaskParams{
			Workspace: "workspace-" + faker.Word(),
			RepoSlug:  "repo-" + faker.Word(),
			PullReqID: 100 + rand.IntN(10000),
			TaskID:    100 + rand.IntN(10000),
		}

		// Act
		result, err := client.GetTask(t.Context(), mockTokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get token")
		assert.ErrorIs(t, errors.Unwrap(err), tokenErr)
	})
}
