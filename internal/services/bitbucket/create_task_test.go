package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreatePullRequestTask(t *testing.T) {
	t.Run("success with all fields", func(t *testing.T) {
		// Arrange
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000 // Convert to a reasonable number
		content := faker.Sentence()
		commentID := faker.RandomUnixTime() // Already int64, no need for conversion
		pending := true

		expectedTask := PullRequestCommentTask{
			PullRequestTask: PullRequestTask{
				Task: Task{
					ID:      faker.RandomUnixTime(), // Already int64, no need for conversion
					State:   "OPEN",
					Content: &TaskContent{Raw: content},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			assert.Equal(t, http.MethodPost, r.Method)
			// Split long line into multiple parts
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks",
				workspace, repoSlug, pullReqID)
			assert.Equal(t, expectedPath, r.URL.Path)

			// Verify auth header
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))

			// Parse request body
			var requestBody CreateTaskPayload
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			assert.NoError(t, err)

			// Verify request payload
			assert.Equal(t, content, requestBody.Content.Raw)
			assert.NotNil(t, requestBody.Comment)
			assert.Equal(t, commentID, requestBody.Comment.ID)
			assert.NotNil(t, requestBody.Pending)
			assert.Equal(t, pending, *requestBody.Pending)

			// Return response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(expectedTask)
			assert.NoError(t, err)
		}))
		defer server.Close()

		// Create client with test-specific logger
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Create token provider
		tokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: "token123",
		}

		// Act
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
			CommentID: commentID,
			Pending:   &pending,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedTask, *result) // Compare entire struct
	})

	t.Run("success with required fields only", func(t *testing.T) {
		// Arrange
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000 // Convert to a reasonable number
		content := faker.Sentence()

		expectedTask := PullRequestCommentTask{
			PullRequestTask: PullRequestTask{
				Task: Task{
					ID:      faker.RandomUnixTime(), // Already int64, no need for conversion
					State:   "OPEN",
					Content: &TaskContent{Raw: content},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			assert.Equal(t, http.MethodPost, r.Method)
			// Split long line into multiple parts
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks",
				workspace, repoSlug, pullReqID)
			assert.Equal(t, expectedPath, r.URL.Path)

			// Verify auth header
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))

			// Parse request body
			var requestBody CreateTaskPayload
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			assert.NoError(t, err)

			// Verify request payload
			assert.Equal(t, content, requestBody.Content.Raw)
			assert.Nil(t, requestBody.Comment)
			assert.Nil(t, requestBody.Pending)

			// Return response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(expectedTask)
			assert.NoError(t, err)
		}))
		defer server.Close()

		// Create client with test-specific logger
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Create token provider
		tokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: "token123",
		}

		// Act
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedTask, *result) // Compare entire struct
	})

	t.Run("api error", func(t *testing.T) {
		// Arrange
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000 // Convert to a reasonable number
		content := faker.Sentence()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify auth header
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))

			// Create a malformed JSON response that will cause a parsing error
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{ "error": "Invalid request", malformed}`)) // This is intentionally malformed
			assert.NoError(t, err)
		}))
		defer server.Close()

		// Create client with test-specific logger
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Create token provider
		tokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: "token123",
		}

		// Act
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("token error", func(t *testing.T) {
		// Arrange
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := int(faker.RandomUnixTime()) % 10000 // Convert to a reasonable number
		content := faker.Sentence()
		tokenErr := errors.New("failed to get token")

		// Create a dummy server just for the client
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			// This should not be called
			assert.Fail(t, "Server should not be called when token provider fails")
		}))
		defer server.Close()

		// Create client with test-specific logger
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Create token provider that returns error
		tokenProvider := &MockTokenProvider{
			Err: tokenErr,
		}

		// Act
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
