package bitbucket

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_MergePR(t *testing.T) {
	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1
		prTitle := "Merged PR " + faker.Sentence()
		prDescription := faker.Paragraph()
		sourceBranch := "feature-" + faker.Word()
		targetBranch := "main-" + faker.Word()
		commitHash := faker.UUIDHyphenated()
		mergeCommitHash := faker.UUIDHyphenated()
		mergeMessage := "Merging " + faker.Sentence()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		updatedOn := time.Now().UTC().Truncate(time.Second)
		updatedOnStr := updatedOn.Format(time.RFC3339)
		createdOnStr := updatedOn.Add(-24 * time.Hour).Format(time.RFC3339)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "POST", r.Method)
			expectedPath := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/merge",
				username, repoSlug, pullRequestID)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"title": "%s",
				"description": "%s",
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
						"name": "%s"
					},
					"commit": {
						"hash": "%s"
					},
					"repository": {
						"full_name": "%s/%s",
						"name": "%s",
						"uuid": "{7708d810-964c-403f-aa6d-4e949280d614}"
					}
				},
				"destination": {
					"branch": {
						"name": "%s"
					},
					"repository": {
						"full_name": "%s/%s",
						"name": "%s",
						"uuid": "{7708d810-964c-403f-aa6d-4e949280d614}"
					}
				},
				"close_source_branch": true,
				"merge_commit": {
					"hash": "%s"
				},
				"comment_count": 0,
				"task_count": 0,
				"created_on": "%s",
				"updated_on": "%s"
			}`, pullRequestID, prTitle, prDescription, sourceBranch, commitHash,
				username, repoSlug, repoSlug, targetBranch, username, repoSlug, repoSlug,
				mergeCommitHash, createdOnStr, updatedOnStr)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
			MergeParameters: &PullRequestMergeParameters{
				CloseSourceBranch: true,
				Message:           mergeMessage,
				MergeStrategy:     "merge_commit",
			},
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, pullRequestID, result.ID)
		assert.Equal(t, prTitle, result.Title)
		assert.Equal(t, "MERGED", result.State)
		assert.Equal(t, sourceBranch, result.Source.Branch.Name)
		assert.Equal(t, targetBranch, result.Destination.Branch.Name)
		assert.True(t, result.CloseSourceBranch)
		assert.Equal(t, mergeCommitHash, result.MergeCommit.Hash)
		assert.Equal(t, updatedOn, result.UpdatedOn.UTC())
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1
		prTitle := "Simple Merge " + faker.Word()
		sourceBranch := "feature-" + faker.Word()
		targetBranch := "main-" + faker.Word()
		mergeCommitHash := faker.UUIDHyphenated()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Return minimal successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"title": "%s",
				"state": "MERGED",
				"source": {
					"branch": {
						"name": "%s"
					}
				},
				"destination": {
					"branch": {
						"name": "%s"
					}
				},
				"merge_commit": {
					"hash": "%s"
				}
			}`, pullRequestID, prTitle, sourceBranch, targetBranch, mergeCommitHash)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
			// No merge parameters, using defaults
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, pullRequestID, result.ID)
		assert.Equal(t, prTitle, result.Title)
		assert.Equal(t, "MERGED", result.State)
		assert.Equal(t, sourceBranch, result.Source.Branch.Name)
		assert.Equal(t, targetBranch, result.Destination.Branch.Name)
		assert.Equal(t, mergeCommitHash, result.MergeCommit.Hash)
	})

	t.Run("handles API error", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

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
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "merge pull request failed")
	})

	t.Run("handles token provider error", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1

		mockTokenProvider := &MockTokenProvider{
			Err: errors.New(faker.Sentence()),
		}

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, "http://example.com")
		client := NewClient(deps)

		// Act
		result, err := client.MergePR(t.Context(), mockTokenProvider, MergePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		expectedError := fmt.Errorf("failed to get token: %w", mockTokenProvider.Err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}
