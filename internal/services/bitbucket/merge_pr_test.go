package bitbucket

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_MergePR(t *testing.T) {
	mockTokenProvider := &MockTokenProvider{}

	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/repositories/test-user/test-repo/pullrequests/1/merge", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"id": 1,
				"title": "Merged PR",
				"description": "Test description",
				"state": "MERGED",
				"author": {
					"account_id": "123456",
					"display_name": "Test User",
					"nickname": "testuser",
					"username": "testuser",
					"uuid": "{58021780-82b6-4517-b153-0ae73ce3e4b4}",
					"type": "user"
				},
				"source": {
					"branch": {
						"name": "feature-branch"
					},
					"commit": {
						"hash": "abcdef123456"
					},
					"repository": {
						"full_name": "test-user/test-repo",
						"name": "test-repo",
						"uuid": "{7708d810-964c-403f-aa6d-4e949280d614}"
					}
				},
				"destination": {
					"branch": {
						"name": "main"
					},
					"repository": {
						"full_name": "test-user/test-repo",
						"name": "test-repo",
						"uuid": "{7708d810-964c-403f-aa6d-4e949280d614}"
					}
				},
				"close_source_branch": true,
				"merge_commit": {
					"hash": "fedcba654321"
				},
				"comment_count": 0,
				"task_count": 0,
				"created_on": "2023-01-01T00:00:00Z",
				"updated_on": "2023-01-02T00:00:00Z"
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDeps(server.URL)
		client := NewClient(deps)

		// Setup token provider
		mockTokenProvider.Token = "test-token"
		mockTokenProvider.Err = nil

		updatedOn, _ := time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")

		// Execute the request
		result, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 1,
			MergeParameters: &PullRequestMergeParameters{
				CloseSourceBranch: true,
				Message:           "Merging feature into main",
				MergeStrategy:     "merge_commit",
			},
		})

		// Verify the result
		require.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Merged PR", result.Title)
		assert.Equal(t, "MERGED", result.State)
		assert.Equal(t, "feature-branch", result.Source.Branch.Name)
		assert.Equal(t, "main", result.Destination.Branch.Name)
		assert.True(t, result.CloseSourceBranch)
		assert.Equal(t, "fedcba654321", result.MergeCommit.Hash)
		assert.Equal(t, updatedOn.UTC(), result.UpdatedOn.UTC())
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return minimal successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{
				"id": 2,
				"title": "Simple Merge",
				"state": "MERGED",
				"source": {
					"branch": {
						"name": "feature-branch"
					}
				},
				"destination": {
					"branch": {
						"name": "main"
					}
				},
				"merge_commit": {
					"hash": "abcdef123456"
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

		// Execute the request with no merge parameters (using defaults)
		result, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 2,
			// No merge parameters, using defaults
		})

		// Verify the result
		require.NoError(t, err)
		assert.Equal(t, 2, result.ID)
		assert.Equal(t, "Simple Merge", result.Title)
		assert.Equal(t, "MERGED", result.State)
		assert.Equal(t, "feature-branch", result.Source.Branch.Name)
		assert.Equal(t, "main", result.Destination.Branch.Name)
		assert.Equal(t, "abcdef123456", result.MergeCommit.Hash)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Setup mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, `{
				"error": {
					"message": "Pull request has conflicts that need to be resolved"
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
		_, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 3,
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "merge pull request failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Create client with mock dependencies
		deps := makeMockDeps("http://example.com")
		client := NewClient(deps)

		// Setup token provider to return an error
		mockTokenProvider.Err = errors.New("token error")

		// Execute the request
		_, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      "test-user",
			RepoSlug:      "test-repo",
			PullRequestID: 1,
		})

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get token")
	})
}
