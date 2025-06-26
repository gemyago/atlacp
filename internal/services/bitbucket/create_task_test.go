package bitbucket

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreatePullRequestTask(t *testing.T) {
	t.Run("success with all fields", func(t *testing.T) {
		// Setup
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := 123
		content := "Task description"
		commentID := int64(456)
		pending := true

		expectedTask := PullRequestCommentTask{
			PullRequestTask: PullRequestTask{
				Task: Task{
					ID:      789,
					State:   "OPEN",
					Content: &TaskContent{Raw: content},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/repositories/"+workspace+"/"+repoSlug+"/pullrequests/123/tasks", r.URL.Path)

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

		// Create client
		client := &Client{
			httpClient: server.Client(),
			baseURL:    server.URL,
		}

		// Create token provider
		tokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: "token123",
		}

		// Execute
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
			CommentID: commentID,
			Pending:   &pending,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedTask.ID, result.ID)
		assert.Equal(t, expectedTask.State, result.State)
		assert.Equal(t, expectedTask.Content.Raw, result.Content.Raw)
	})

	t.Run("success with required fields only", func(t *testing.T) {
		// Setup
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := 123
		content := "Task description"

		expectedTask := PullRequestCommentTask{
			PullRequestTask: PullRequestTask{
				Task: Task{
					ID:      789,
					State:   "OPEN",
					Content: &TaskContent{Raw: content},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/repositories/"+workspace+"/"+repoSlug+"/pullrequests/123/tasks", r.URL.Path)

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

		// Create client
		client := &Client{
			httpClient: server.Client(),
			baseURL:    server.URL,
		}

		// Create token provider
		tokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: "token123",
		}

		// Execute
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedTask.ID, result.ID)
		assert.Equal(t, expectedTask.State, result.State)
		assert.Equal(t, expectedTask.Content.Raw, result.Content.Raw)
	})

	t.Run("api error", func(t *testing.T) {
		// Setup
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := 123
		content := "Task description"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Create a malformed JSON response that will cause a parsing error
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(`{ "error": "Invalid request", malformed}`)) // This is intentionally malformed
			assert.NoError(t, err)
		}))
		defer server.Close()

		// Create a client with timeout to ensure the test fails quickly
		client := &Client{
			httpClient: &http.Client{
				Timeout: 1 * time.Second,
			},
			baseURL: server.URL,
		}

		// Create token provider
		tokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: "token123",
		}

		// Execute
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Verify
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("token error", func(t *testing.T) {
		// Setup
		workspace := faker.Username()
		repoSlug := faker.Username()
		pullReqID := 123
		content := "Task description"
		tokenErr := errors.New("failed to get token")

		// Create client
		client := &Client{
			httpClient: &http.Client{},
			baseURL:    "https://api.bitbucket.org/2.0",
		}

		// Create token provider that returns error
		tokenProvider := &MockTokenProvider{
			Err: tokenErr,
		}

		// Execute
		params := CreatePullRequestTaskParams{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			PullReqID: pullReqID,
			Content:   content,
		}
		result, err := client.CreatePullRequestTask(t.Context(), tokenProvider, params)

		// Verify
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
