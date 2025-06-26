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

func TestClient_UpdatePR(t *testing.T) {
	t.Run("success with all parameters and fields", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1
		updatedTitle := "Updated PR " + faker.Sentence()
		updatedDescription := faker.Paragraph()
		destinationBranch := "main-" + faker.Word()

		mockTokenProvider := &MockTokenProvider{
			TokenType:  "Bearer",
			TokenValue: faker.UUIDHyphenated(),
		}

		updatedOn := time.Now().UTC().Truncate(time.Second)
		updatedOnStr := updatedOn.Format(time.RFC3339)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request details
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", username, repoSlug, pullRequestID), r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer "+mockTokenProvider.TokenValue, r.Header.Get("Authorization"))

			// Return complete successful response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{
				"id": %d,
				"title": "%s",
				"description": "%s",
				"state": "OPEN",
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
				"comment_count": 0,
				"task_count": 0,
				"created_on": "2023-01-01T00:00:00Z",
				"updated_on": "%s"
			}`, pullRequestID, updatedTitle, updatedDescription, username, repoSlug,
				repoSlug, destinationBranch, username, repoSlug, repoSlug, updatedOnStr)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.UpdatePR(t.Context(), mockTokenProvider, UpdatePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
			Request: &PullRequest{
				Title:       updatedTitle,
				Description: updatedDescription,
				Destination: &PullRequestDestination{
					Branch: PullRequestBranch{
						Name: destinationBranch,
					},
				},
				CloseSourceBranch: true,
			},
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, pullRequestID, result.ID)
		assert.Equal(t, updatedTitle, result.Title)
		assert.Equal(t, updatedDescription, result.Description)
		assert.Equal(t, "OPEN", result.State)
		assert.Equal(t, "feature-branch", result.Source.Branch.Name)
		assert.Equal(t, destinationBranch, result.Destination.Branch.Name)
		assert.True(t, result.CloseSourceBranch)
		assert.Equal(t, updatedOn, result.UpdatedOn.UTC())
	})

	t.Run("success with required parameters only", func(t *testing.T) {
		// Arrange
		username := "test-user-" + faker.Word()
		repoSlug := "test-repo-" + faker.Word()
		pullRequestID := rand.Intn(1000) + 1
		updatedTitle := "Updated Title " + faker.Word()

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
				"state": "OPEN",
				"source": {
					"branch": {
						"name": "feature-branch"
					}
				}
			}`, pullRequestID, updatedTitle)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.UpdatePR(t.Context(), mockTokenProvider, UpdatePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
			Request: &PullRequest{
				Title: updatedTitle,
			},
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, pullRequestID, result.ID)
		assert.Equal(t, updatedTitle, result.Title)
		assert.Equal(t, "OPEN", result.State)
		assert.Equal(t, "feature-branch", result.Source.Branch.Name)
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{
				"error": {
					"message": "Invalid update request"
				}
			}`)
		}))
		defer server.Close()

		// Create client with mock dependencies
		deps := makeMockDepsWithTestName(t, server.URL)
		client := NewClient(deps)

		// Act
		result, err := client.UpdatePR(t.Context(), mockTokenProvider, UpdatePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
			Request:       &PullRequest{
				// Empty update request
			},
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, "update pull request failed")
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
		result, err := client.UpdatePR(t.Context(), mockTokenProvider, UpdatePRParams{
			Username:      username,
			RepoSlug:      repoSlug,
			PullRequestID: pullRequestID,
			Request: &PullRequest{
				Title: faker.Sentence(),
			},
		})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		expectedError := fmt.Errorf("failed to get token: %w", mockTokenProvider.Err)
		assert.Equal(t, expectedError.Error(), err.Error())
	})
}
