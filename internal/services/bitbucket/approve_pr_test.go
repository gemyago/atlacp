package bitbucket

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ApprovePR(t *testing.T) {
	mockTokenProvider := &MockTokenProvider{}

	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/repositories/test-user/test-repo/pullrequests/1/approve", r.URL.Path)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"user": {
					"account_id": "123456",
					"display_name": "Test User",
					"nickname": "testuser",
					"username": "testuser",
					"uuid": "{58021780-82b6-4517-b153-0ae73ce3e4b4}",
					"type": "user"
				},
				"role": "REVIEWER",
				"approved": true,
				"state": "approved",
				"type": "participant"
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request
		result, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 1,
		})

		// Verify the result
		require.NoError(t, err)
		assert.Equal(t, "Test User", result.User.DisplayName)
		assert.Equal(t, "testuser", result.User.Username)
		assert.Equal(t, "REVIEWER", result.Role)
		assert.True(t, result.Approved)
		assert.Equal(t, "approved", result.State)
		assert.Equal(t, "participant", result.Type)
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return minimal successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"user": {
					"username": "testuser",
					"type": "user"
				},
				"approved": true,
				"type": "participant"
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request
		result, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 1,
		})

		// Verify the result
		require.NoError(t, err)
		assert.Equal(t, "testuser", result.User.Username)
		assert.True(t, result.Approved)
		assert.Equal(t, "participant", result.Type)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{
				"error": {
					"message": "Cannot approve your own pull request"
				}
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		// Execute the request
		_, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 1,
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "approve pull request failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Create client with mock dependencies
		deps := makeMockDeps("http://example.com")
		client := NewClient(deps)

		// Setup token provider to return an error
		mockTokenProvider.Err = errors.New("token error")

		// Execute the request
		_, err := client.ApprovePR(t.Context(), mockTokenProvider, ApprovePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 1,
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
